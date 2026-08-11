import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";

// blog.nagutabby.uk/api.nagutabby.uk用のRoute 53ホストゾーンとACM証明書はmigrate/16で追加する。
export class DnsStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
  }
}
