import { useEffect, useState } from 'react';
import { Button, Drawer } from 'antd';
import { RightOutlined } from '@ant-design/icons';

const CalcPrice = async (data: any[]) => {
  const response = await fetch(`calc-price`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Failed to calculate price');
  const price = await response.json();
  console.log(price);
  return price;
};

const PutData = (data: any[]) => {
  // eg. data = ['X-1', 'W-1', 'Y-3', 'Z-2']
  // let data looks like this: X-1 + W-1 + Y-3 + Z-2
  const returnData = data.join(' + ');
  console.log(returnData);
  return returnData;
}


const Result = ({ data, render }: { data: any[], render: () => boolean }) => {
  const [price, setPrice] = useState<number | null>(null);
  const [show, setShow] = useState(false);
  useEffect(() => {
    CalcPrice(data).then(setPrice);
  }, [data]);
  return (
    <>
    <Button type="primary" shape='round' icon={<RightOutlined />} onClick={()=>setShow(render)}>Next</Button>
    <Drawer title="Selected Components" placement="right" closable={true} onClose={()=>setShow(false)} open={show}>
      <p>{PutData(data)}</p>
      <h2>Total: ${price}</h2>
    </Drawer>
    </>
  );
};

export default Result;