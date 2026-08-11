import * as path from "node:path";
import { Duration, RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { Architecture, FunctionUrlAuthType, Runtime } from "aws-cdk-lib/aws-lambda";
import { Secret } from "aws-cdk-lib/aws-secretsmanager";
import { ITable } from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";

export interface ApiStackProps extends StackProps {
  readonly followerTable: ITable;
  readonly relayConnectionTable: ITable;
  // お問い合わせフォームの送信元/BCC先。実際のアドレスはリポジトリに書かず、
  // `cdk deploy -c fromAddress=... -c bccAddress=...`のようにcontextで渡す。
  readonly fromAddress: string;
  readonly bccAddress: string;
}

// フォロワーへの署名付き配送がLambdaのデフォルトタイムアウト(3秒)を超えうるため、
// 30秒に延長する。将来フォロワーが増えてfan-outが遅くなった場合はSQSでの非同期化を検討する。
const HANDLER_TIMEOUT = Duration.seconds(30);

export class ApiStack extends Stack {
  readonly functionUrl: string;

  constructor(scope: Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

    // ACTOR_PRIVATE_KEY_PEM/ACTOR_PUBLIC_KEY_PEM/EMAIL_API_TOKENの実際の値は
    // `aws secretsmanager put-secret-value`でRailwayの既存値から手動投入する
    // (このリポジトリにはコミットしない)。
    const actorKeys = new Secret(this, "ActorKeysSecret", {
      description: "ActivityPub actorのHTTP Signature用鍵ペア(PEM)",
      removalPolicy: RemovalPolicy.RETAIN,
    });
    const emailApiToken = new Secret(this, "EmailApiTokenSecret", {
      description: "お問い合わせフォームのメール送信APIトークン",
      removalPolicy: RemovalPolicy.RETAIN,
    });

    const handler = new GoFunction(this, "BackendFunction", {
      entry: path.join(__dirname, "..", "..", "backend", "cmd", "lambda"),
      architecture: Architecture.ARM_64,
      runtime: Runtime.PROVIDED_AL2023,
      timeout: HANDLER_TIMEOUT,
      environment: {
        SITE_BASE_URL: "https://blog.nagutabby.uk",
        WEB_BASE_URL: "https://blog.nagutabby.uk",
        FOLLOWER_TABLE_NAME: props.followerTable.tableName,
        RELAY_CONNECTION_TABLE_NAME: props.relayConnectionTable.tableName,
        ACTOR_KEYS_SECRET_ARN: actorKeys.secretArn,
        EMAIL_API_TOKEN_SECRET_ARN: emailApiToken.secretArn,
        FROM_ADDRESS: props.fromAddress,
        BCC_ADDRESS: props.bccAddress,
      },
    });

    props.followerTable.grantReadWriteData(handler);
    props.relayConnectionTable.grantReadWriteData(handler);
    actorKeys.grantRead(handler);
    emailApiToken.grantRead(handler);

    // migrate/14でCloudFrontをFunction URLの手前に置くまでの検証用。
    // 認証なしで公開する(webfinger/actor/inbox等は元々認証なしの公開エンドポイントのため)。
    const fnUrl = handler.addFunctionUrl({
      authType: FunctionUrlAuthType.NONE,
    });
    this.functionUrl = fnUrl.url;
  }
}
