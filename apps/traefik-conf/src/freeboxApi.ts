import * as freebox from "@mgarnier11/freebox";
import { config } from "./utils/config.js";

const findRedirect = (redirectionName: string, redirectionData: any) => {
  const redirs = redirectionData.result;

  const ovhRedir = redirs.find((redir: any) => redir.comment.toLowerCase() === redirectionName.toLowerCase());

  return ovhRedir;
};

const getFbx = async () => {
  const fbx = new freebox.Freebox({
    app_token: config.fbxAppToken,
    app_id: config.fbxAppId,
    api_domain: config.fbxApiDomain,
    https_port: config.fbxHttpsPort,
    api_base_url: config.fbxApiBaseUrl,
    api_version: config.fbxApiVersion,
  });

  await fbx.login();

  return fbx;
};

export const getRedirection = async (redirectionName: string) => {
  const fbx = await getFbx();

  const { data: redirectionData } = await fbx.request({
    method: "GET",
    url: "fw/redir/",
  });

  const redirection = findRedirect(redirectionName, redirectionData);

  return redirection;
};

export const changeRedirection = async (redirectionName: string, newIp: string) => {
  const fbx = await getFbx();

  const { data: redirectionData } = await fbx.request({
    method: "GET",
    url: "fw/redir/",
  });

  const redirection = findRedirect(redirectionName, redirectionData);

  const { data: resultData } = await fbx.request({
    method: "PUT",
    url: `fw/redir/${redirection.id}/`,
    data: {
      lan_ip: newIp,
    },
  });

  await fbx.logout();

  if (resultData.success === true) {
    return true;
  } else {
    return false;
  }
};
