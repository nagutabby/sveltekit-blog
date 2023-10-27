type Brand = {
  readonly brand: string;
  readonly version: string;
};

type NavigatorUAData = {
  readonly brands: Brand[];
  readonly mobile: boolean;
  readonly platform: string;
};

interface Navigator {
  userAgentData: NavigatorUAData;
  msMaxTouchPoints: number;
}

interface CSSStyleDeclaration {
  msOverflowStyle: string;
}
