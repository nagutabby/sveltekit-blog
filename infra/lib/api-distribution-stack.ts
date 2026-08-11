import { Stack, StackProps } from "aws-cdk-lib";
import { AllowedMethods, CachePolicy, Distribution, OriginRequestPolicy, ViewerProtocolPolicy } from "aws-cdk-lib/aws-cloudfront";
import { FunctionUrlOrigin } from "aws-cdk-lib/aws-cloudfront-origins";
import { FunctionUrl } from "aws-cdk-lib/aws-lambda";
import { ICertificate } from "aws-cdk-lib/aws-certificatemanager";
import { ARecord, AaaaRecord, IHostedZone, RecordTarget } from "aws-cdk-lib/aws-route53";
import { CloudFrontTarget } from "aws-cdk-lib/aws-route53-targets";
import { Construct } from "constructs";

const DOMAIN_NAME = "api.nagutabby.uk";

export interface ApiDistributionStackProps extends StackProps {
  readonly functionUrl: FunctionUrl;
  readonly hostedZone: IHostedZone;
  readonly certificate: ICertificate;
}

// api.nagutabby.uk全体が同じLambda(ApiStack)への単純なパススルーであり、
// S3等の他オリジンを持たないため、defaultBehaviorだけで足りる
// (blog.nagutabby.ukのSiteStackのようにパスごとの振り分けは不要)。
export class ApiDistributionStack extends Stack {
  constructor(scope: Construct, id: string, props: ApiDistributionStackProps) {
    super(scope, id, props);

    const backendOrigin = FunctionUrlOrigin.withOriginAccessControl(props.functionUrl);

    const distribution = new Distribution(this, "ApiDistribution", {
      domainNames: [DOMAIN_NAME],
      certificate: props.certificate,
      defaultBehavior: {
        origin: backendOrigin,
        viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        cachePolicy: CachePolicy.CACHING_DISABLED,
        allowedMethods: AllowedMethods.ALLOW_ALL,
        originRequestPolicy: OriginRequestPolicy.ALL_VIEWER_EXCEPT_HOST_HEADER,
      },
    });

    const aliasTarget = RecordTarget.fromAlias(new CloudFrontTarget(distribution));
    new ARecord(this, "ApiAliasRecordV4", { zone: props.hostedZone, target: aliasTarget });
    new AaaaRecord(this, "ApiAliasRecordV6", { zone: props.hostedZone, target: aliasTarget });
  }
}
