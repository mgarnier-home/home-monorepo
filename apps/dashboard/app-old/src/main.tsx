import './index.css';

import { logger } from 'logger';
import ReactDOM from 'react-dom/client';

import App from './components/app';

logger.setAppName('dashboard-client');

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  // <React.StrictMode>
  <App />
  // </React.StrictMode>
);
