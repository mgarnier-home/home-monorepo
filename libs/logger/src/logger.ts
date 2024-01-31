import { Color, LogLevel } from './interfaces';

const envLogLevel = parseInt(process.env.LOG_LEVEL);
const logLevel = isNaN(envLogLevel) ? LogLevel.INFO : envLogLevel;

interface ILogger {
  // log(color: Color, ...args: any[]): void;
  info(...args: any[]): void;
  debug(...args: any[]): void;
  warn(...args: any[]): void;
  error(...args: any[]): void;
  verbose(...args: any[]): void;
}

class Logger implements ILogger {
  public colored: ColoredLogger = new ColoredLogger(this);
  private appName: string;

  private static getColor(def: Color, arg: any): Color {
    return typeof arg === 'string' && arg.startsWith('\x1b') ? (arg as Color) : def;
  }

  private log(color: Color, level: LogLevel, ...args: any[]) {
    const date = new Date().toISOString();

    console.log(color, `[${date}] [${level}] [${this.appName}] `, ...args, '\x1b[0m');
  }

  private errorLog(color: Color, level: LogLevel, ...args: any[]) {
    const date = new Date().toISOString();

    console.error(color, `[${date}] [${level}] [${this.appName}] `, ...args, '\x1b[0m');
  }

  setAppName(appName: string) {
    this.appName = appName;
  }

  info(...args: any[]) {
    if (logLevel <= LogLevel.INFO) {
      this.log(Logger.getColor(Color.DEFAULT, args[0]), LogLevel.INFO, ...args);
    }
  }

  debug(...args: any[]) {
    if (logLevel <= LogLevel.DEBUG) {
      this.log(Logger.getColor(Color.BLUE, args[0]), LogLevel.DEBUG, ...args);
    }
  }

  warn(...args: any[]) {
    if (logLevel <= LogLevel.WARN) {
      this.log(Logger.getColor(Color.YELLOW, args[0]), LogLevel.WARN, ...args);
    }
  }

  error(...args: any[]) {
    if (logLevel <= LogLevel.ERROR) {
      this.errorLog(Logger.getColor(Color.RED, args[0]), LogLevel.ERROR, ...args);
    }
  }

  verbose(...args: any[]) {
    if (logLevel <= LogLevel.VERBOSE) {
      this.log(Logger.getColor(Color.GRAY, args[0]), LogLevel.VERBOSE, ...args);
    }
  }
}

class ColoredLogger implements ILogger {
  constructor(private logger: Logger) {}

  info(color: Color, ...args: any[]) {
    this.logger.info(color, ...args);
  }

  debug(color: Color, ...args: any[]) {
    this.logger.debug(color, ...args);
  }

  warn(color: Color, ...args: any[]) {
    this.logger.warn(color, ...args);
  }

  error(color: Color, ...args: any[]) {
    this.logger.error(color, ...args);
  }

  verbose(color: Color, ...args: any[]) {
    this.logger.verbose(color, ...args);
  }
}

export const logger = new Logger();
