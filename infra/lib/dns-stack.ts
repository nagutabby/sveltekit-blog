import { CfnOutput, Fn, Stack, StackProps } from "aws-cdk-lib";
import { Certificate, CertificateValidation, ICertificate } from "aws-cdk-lib/aws-certificatemanager";
import { HostedZone, IHostedZone } from "aws-cdk-lib/aws-route53";
import { Construct } from "constructs";

const BLOG_DOMAIN = "blog.nagutabby.uk";
const API_DOMAIN = "api.nagutabby.uk";

export interface DnsStackZone {
  readonly hostedZone: IHostedZone;
  readonly certificate: ICertificate;
}

// nagutabby.ukゾーン自体はCloudflareに残る。blog/apiサブドメインだけをNSレコードで
// Route 53に委任する(ゾーンapexやメール関連レコードには触れない、フルゾーン移行より
// リスクが小さい)。ここでは新規ゾーン+証明書の作成のみを行い、実際のNS委任
// (Cloudflare側にNSレコードを追加する)は手動オペレーションとして別に行う。
//
// 証明書のDNS検証とCloudFrontへのエイリアスレコード作成には循環依存の関係がある
// (Distributionは作成時にcertificate/domainNamesが必要だが、エイリアスレコードは
// Distributionの作成後でないと作れない)。そのため、ゾーン+証明書の作成はここ
// (DnsStack)に閉じ、エイリアスレコード自体はSiteStack/ApiDistributionStackが
// 自分のDistributionを作った後に自分で作成する(このスタックからは作らない)。
export class DnsStack extends Stack {
  readonly blog: DnsStackZone;
  readonly api: DnsStackZone;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    this.blog = this.createZone("Blog", BLOG_DOMAIN);
    this.api = this.createZone("Api", API_DOMAIN);
  }

  private createZone(idPrefix: string, domainName: string): DnsStackZone {
    const hostedZone = new HostedZone(this, `${idPrefix}Zone`, { zoneName: domainName });
    const certificate = new Certificate(this, `${idPrefix}Certificate`, {
      domainName,
      validation: CertificateValidation.fromDns(hostedZone),
    });

    new CfnOutput(this, `${idPrefix}NameServers`, {
      description: `Cloudflareの${domainName.split(".")[0]}サブドメインにNSレコードとして追加する値`,
      // hostedZoneNameServersはCloudFormationのトークンリストなので、素の
      // Array.prototype.join()ではなくFn.joinを使う必要がある
      // (素のjoinはCDK合成時に"encoded list token in scalar context"エラーになる)。
      value: Fn.join(", ", hostedZone.hostedZoneNameServers ?? []),
    });

    return { hostedZone, certificate };
  }
}
