export namespace ApiInterfaces {
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

  export namespace StatusChecks {
    export type RequestData = { id: string; url: string };

    export type Request = {
      statusChecks: RequestData[];
    };

    export type ResponseData = { id: string; code: number; duration: number };

    export type Response = {
      statusChecks: ResponseData[];
    };
  }
}
