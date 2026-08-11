import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";

// blog.nagutabby.uk向けのS3+CloudFront(SSGサイト+Federationパスのビヘイビア分岐)はmigrate/14で追加する。
export class SiteStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
  }
}
