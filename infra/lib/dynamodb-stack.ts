import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";

// Follower/RelayConnectionテーブルはmigrate/12で追加する。
export class DynamoDbStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
  }
}
