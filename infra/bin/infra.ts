#!/usr/bin/env node
import { App } from "aws-cdk-lib";
import { DynamoDbStack } from "../lib/dynamodb-stack";
import { ApiStack } from "../lib/api-stack";
import { SiteStack } from "../lib/site-stack";
import { ApiDistributionStack } from "../lib/api-distribution-stack";
import { DnsStack } from "../lib/dns-stack";

// CloudFrontに証明書をアタッチするにはus-east-1のACM証明書が必要なため、
// リージョン分割の複雑さを避けてスタック全体をus-east-1に統一する。
const env = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: "us-east-1",
};

const app = new App();

// `cdk deploy -c fromAddress=... -c bccAddress=...`で渡す。cdk synthだけなら
// 未指定でも空文字列のまま通る(実際のメールアドレスをリポジトリに書かないため)。
const fromAddress = (app.node.tryGetContext("fromAddress") as string | undefined) ?? "";
const bccAddress = (app.node.tryGetContext("bccAddress") as string | undefined) ?? "";

const dynamoDbStack = new DynamoDbStack(app, "SveltekitBlogDynamoDbStack", { env });
new ApiStack(app, "SveltekitBlogApiStack", {
  env,
  followerTable: dynamoDbStack.followerTable,
  relayConnectionTable: dynamoDbStack.relayConnectionTable,
  fromAddress,
  bccAddress,
});
new SiteStack(app, "SveltekitBlogSiteStack", { env });
new ApiDistributionStack(app, "SveltekitBlogApiDistributionStack", { env });
new DnsStack(app, "SveltekitBlogDnsStack", { env });
