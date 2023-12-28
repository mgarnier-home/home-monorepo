export namespace ApiInterfaces {
  export namespace StatusCheck {
    export type Request = {
      url: string;
      method: string;
    };

    export type Response = {
      code: number;
      duration: number;
    };
  }

  export namespace PingHost {
    export type Request = {
      ip: string;
    };

    export type Response = {
      ping: boolean;
      duration: number;
      ms: number;
    };
  }

  export namespace MakeRequest {
    export type Request = {
      url: string;
      method: string;
      body?: string;
    };

    export type Response<Data> = {
      code: number;
      duration: number;
      data?: Data;
    };
  }
}
