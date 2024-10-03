import '../App.css';
import { Space, Button } from 'antd';
import viteLogo from '/vite.svg'
import reactLogo from '../assets/react.svg';
import { Link } from 'react-router-dom';
import Counter from '../features/counter/Counter'
function Home() {
  return (
    <Space direction="vertical">
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <Button type="primary"><Link to='select'>Go Select</Link></Button>
      <Counter />
      <p>Edit <code>src/App.tsx</code> and save to test HMR</p>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </Space>
  );
}

export default Home;