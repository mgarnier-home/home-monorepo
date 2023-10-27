import { useEffect, useState } from "react";

import AppInterfaces from "@shared/interfaces/appInterfaces";

import { Api } from "../utils/api";
import { ConfigContext } from "../utils/configContext";
import { socket } from "../utils/socket";
import Utils from "../utils/utils";
import HomeComponent from "./home/home";

function App() {
  const [config, setConfig] = useState({} as AppInterfaces.AppConfig);

  useEffect(() => {
    console.log("App mounted");

    socket.connect();

    const getConfig = async () => {
      const config = await Api.getConfig();

      setConfig(config);

      document.title = config.pageConfig.pageTitle || "Dashboard";
    };

    getConfig();

    return () => {
      console.log("App unmounted");
      socket.dispose();
    };
  }, []);

  return (
    <div>
      <ConfigContext.Provider value={config}>
        <HomeComponent hosts={config?.hosts || []} />
      </ConfigContext.Provider>
    </div>
  );
}

export default App;
