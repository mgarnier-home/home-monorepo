import './test.css';

import ReactDOM from 'react-dom/client';
import { t2 } from 'utils';

import { t3 } from '../../shared/utils';

// App Component
const App = () => (<div>
  <h1 className={'test'}>Hello, ESBUILD!</h1>
  <Panel />
  <Panel />
</div>)

// Panel Component
const Panel = () => <h2>I'm a Panel</h2>

// Mount component 
const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(<App />);

t2();

t2();

t3();