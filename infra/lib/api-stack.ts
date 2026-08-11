import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";

// Go LambdaとFunction URL、Secrets Managerの定義はmigrate/13で追加する。
export class ApiStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
  }
}
