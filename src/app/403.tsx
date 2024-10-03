import { Empty } from 'antd';

export default function forbidden() {
    return (
      <>
        <Empty description="Permission denied" />
        <a href="/#/">Go Home</a>
      </>
    )
}