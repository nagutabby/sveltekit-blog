import { RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { AttributeType, BillingMode, Table } from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";

// フォロワー・リレー接続ともに書き込み頻度が低く、無料利用枠(プロビジョンド25RCU/25WCU)を
// 恒久的に使い続けるためBillingMode.PROVISIONEDを使う。オンデマンドは無料枠の対象外。
const READ_CAPACITY = 1;
const WRITE_CAPACITY = 1;
const MAX_CAPACITY = 5;

export class DynamoDbStack extends Stack {
  readonly followerTable: Table;
  readonly relayConnectionTable: Table;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    this.followerTable = new Table(this, "FollowerTable", {
      tableName: "Follower",
      partitionKey: { name: "ActorId", type: AttributeType.STRING },
      billingMode: BillingMode.PROVISIONED,
      readCapacity: READ_CAPACITY,
      writeCapacity: WRITE_CAPACITY,
      removalPolicy: RemovalPolicy.RETAIN,
    });

    this.relayConnectionTable = new Table(this, "RelayConnectionTable", {
      tableName: "RelayConnection",
      partitionKey: { name: "ActorId", type: AttributeType.STRING },
      billingMode: BillingMode.PROVISIONED,
      readCapacity: READ_CAPACITY,
      writeCapacity: WRITE_CAPACITY,
      removalPolicy: RemovalPolicy.RETAIN,
    });

    for (const table of [this.followerTable, this.relayConnectionTable]) {
      table.autoScaleReadCapacity({ minCapacity: READ_CAPACITY, maxCapacity: MAX_CAPACITY }).scaleOnUtilization({
        targetUtilizationPercent: 70,
      });
      table.autoScaleWriteCapacity({ minCapacity: WRITE_CAPACITY, maxCapacity: MAX_CAPACITY }).scaleOnUtilization({
        targetUtilizationPercent: 70,
      });
    }
  }
}
