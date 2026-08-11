import { Stack, StackProps } from "aws-cdk-lib";
import { AllowedMethods, CachePolicy, Distribution, OriginRequestPolicy, ViewerProtocolPolicy } from "aws-cdk-lib/aws-cloudfront";
import { FunctionUrlOrigin } from "aws-cdk-lib/aws-cloudfront-origins";
import { FunctionUrl } from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";

export interface ApiDistributionStackProps extends StackProps {
  readonly functionUrl: FunctionUrl;
}

// api.nagutabby.uk全体が同じLambda(ApiStack)への単純なパススルーであり、
// S3等の他オリジンを持たないため、defaultBehaviorだけで足りる
// (blog.nagutabby.ukのSiteStackのようにパスごとの振り分けは不要)。
export class ApiDistributionStack extends Stack {
  constructor(scope: Construct, id: string, props: ApiDistributionStackProps) {
    super(scope, id, props);

    const backendOrigin = FunctionUrlOrigin.withOriginAccessControl(props.functionUrl);

    new Distribution(this, "ApiDistribution", {
      defaultBehavior: {
        origin: backendOrigin,
        viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        cachePolicy: CachePolicy.CACHING_DISABLED,
        allowedMethods: AllowedMethods.ALLOW_ALL,
        originRequestPolicy: OriginRequestPolicy.ALL_VIEWER_EXCEPT_HOST_HEADER,
      },
    });
  }
}
