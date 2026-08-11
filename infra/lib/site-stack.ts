import * as path from "node:path";
import { RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { Bucket, BlockPublicAccess } from "aws-cdk-lib/aws-s3";
import { BucketDeployment, Source, CacheControl } from "aws-cdk-lib/aws-s3-deployment";
import {
  AllowedMethods,
  CachePolicy,
  Distribution,
  OriginRequestPolicy,
  ViewerProtocolPolicy,
} from "aws-cdk-lib/aws-cloudfront";
import { S3BucketOrigin, FunctionUrlOrigin } from "aws-cdk-lib/aws-cloudfront-origins";
import { FunctionUrl } from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";

export interface SiteStackProps extends StackProps {
  readonly functionUrl: FunctionUrl;
}

// 記事更新をブラウザ/CDNキャッシュ経由で遅延させない([[ssg_unify_rendering]]の
// 既存要件、以前はCloudflare WorkerがCache-Control: no-storeを強制していた)。
const NO_CACHE_POLICY = CachePolicy.CACHING_DISABLED;

// Federation関連パス。Actor識別子(SITE_BASE_URL)がblog.nagutabby.ukに紐づいて
// いるため、これらはblog.nagutabby.uk上でLambdaへ振り分ける必要がある
// (api.nagutabby.uk側には移動できない)。既存のCloudflare Worker
// (web/src/lib/workers/router.ts)のBACKEND_PATH_PREFIXESのうち、
// Federation以外(お問い合わせ)はmigrate/16のapi.nagutabby.ukに移す。
const FEDERATION_PATH_PATTERNS = ["/actor", "/actor/*", "/.well-known/*", "/api/articles/*"];

export class SiteStack extends Stack {
  constructor(scope: Construct, id: string, props: SiteStackProps) {
    super(scope, id, props);

    // ビルド成果物(web/build)は再生成可能なだけなので、DynamoDBテーブルとは異なり
    // DESTROYでよい。
    const bucket = new Bucket(this, "SiteBucket", {
      blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
      removalPolicy: RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
    });

    const s3Origin = S3BucketOrigin.withOriginAccessControl(bucket);
    const backendOrigin = FunctionUrlOrigin.withOriginAccessControl(props.functionUrl);

    const backendBehavior = {
      origin: backendOrigin,
      cachePolicy: NO_CACHE_POLICY,
      // Federationは(request-target)以外にDigest/Signature等の全ヘッダー、
      // InboxはPOSTボディが必須。
      allowedMethods: AllowedMethods.ALLOW_ALL,
      originRequestPolicy: OriginRequestPolicy.ALL_VIEWER_EXCEPT_HOST_HEADER,
    };

    const distribution = new Distribution(this, "SiteDistribution", {
      defaultRootObject: "index.html",
      defaultBehavior: {
        origin: s3Origin,
        viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        cachePolicy: NO_CACHE_POLICY,
      },
      additionalBehaviors: Object.fromEntries(
        FEDERATION_PATH_PATTERNS.map((pattern) => [pattern, backendBehavior]),
      ),
    });

    new BucketDeployment(this, "SiteDeployment", {
      sources: [Source.asset(path.join(__dirname, "..", "..", "web", "build"))],
      destinationBucket: bucket,
      distribution,
      distributionPaths: ["/*"],
      cacheControl: [CacheControl.fromString("no-store")],
      prune: true,
    });
  }
}
