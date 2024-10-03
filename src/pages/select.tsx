import Result from '../app/result';
import React, { useState, useEffect } from 'react';
import { Flex, Radio, Popover, message } from 'antd';

interface Option {
  name: string;
  desc: string;
  price: string;
}

const App: React.FC = () => {
  const [msg, contextHolder] = message.useMessage();
  const [prods, setProds] = useState([]);
  const [currentProd, setCurrentProd] = useState('');
  const [prodDesc, setProdDesc] = useState('');
  const [comp, setComp] = useState([]);
  const [rcomp, setRcomp] = useState<string[]>([]);
  const [option, setOption] = useState<Record<string, Option[]>>({});
  const [selOpt, setSelOpt] = useState<any[]>([]);
  const [selected, setSelected] = useState<string[]>([]);

  const ListProd = async () => {
    try {
      const response = await fetch(`/list-products?t=Product&&col=name`, {
        method: 'GET',
      });
      if (!response.ok) {
        throw new Error('Failed to list Products');
      }
      const data = await response.json();
      setProds(data);
    } catch (error) {
      console.error('Error listing products:', error);
    }
  };

  useEffect(() => {
    ListProd();
  } , []);
  const ProdChange = async (e: any) => {
    setComp([]);
    setRcomp([]);
    setSelOpt([]);
    setSelected([]);
    setProdDesc('');
    setCurrentProd(e.target.value);
    try {
      const response = await fetch(`get-product?name=${e.target.value}`, {
        method: 'GET',
      });
      if (!response.ok) throw new Error('Failed to fetch product');
      const product = await response.json();
      setProdDesc(product.desc);
      setComp(product.comp);
      setRcomp([product.comp[0]]);
      listOptions(product.comp[0]);
    } catch (error) {
      console.error(error);
    }
  };
  const listOptions = async (t: string) => {
    // console.log(t);
    try {
      const response = await fetch(`get-comp-name?t=${t}`, {
        method: 'GET',
      });
      if (!response.ok) throw new Error('Failed to fetch options');
      const options = await response.json();
      const opt = Object.keys(options).map(key => ({
        name: options[key].name,
        desc: options[key].desc,
        price: options[key].price,
      }));
      // console.log(opt); //debug
      setOption((prevOption) => ({ ...prevOption, [t]: opt }));
    } catch (error) {
      console.error(error);
    }
  }
  const onSelect = async (e: any) => {
    // find out the index of the current component type and +1
    const index = comp.findIndex((c) => (e.target.value).includes(c)) + 1;
    if (rcomp.includes(comp[index])) {
      console.log('Already in rcomp:', rcomp.includes(comp[index]));
      // clear options for the next component
      setOption((prevOption) => ({ ...prevOption, [comp[index]]: [] }));
      setRcomp((prevRcomp) => prevRcomp.slice(0, index));
      setSelected((prevSelected) => prevSelected.slice(0, index-1));
      setSelOpt([]);
    }
    // append the selected component to selected
    setSelected((prevSelected) => [...prevSelected, e.target.value]);
    if (index === comp.length) return;
    try {
      const response = await fetch(`pair-components?prod=${currentProd}&&name=${e.target.value}&&next=${comp[index]}`, {
        method: 'GET',
      });
      if (!response.ok) throw new Error(await response.text());
      
      const data = await response.json();
      // console.log("paired data", data); //debug
      // {W: [{desc: "desc_W1", name: "W-1", price: "65"}, {desc: "desc_W2", name: "W-2", price: "35"}]}
      // set object index to the current component type
      const tdOpt = Object.keys(data).map(key => (
        data[key].map((item: any) => ({
          name: item.name,
          desc: item.desc,
          price: item.price,
        }))
      )).flat();
      setSelOpt(tdOpt);
    } catch (error) {
      console.error('Failed to pair components',error);
    }
  }
  useEffect(() => {
    // console.log('rcomp:', rcomp);  //debug
    // console.log('selOpt:', selOpt);  //debug
    for (let i = 0; i < selOpt.length; i++) {
      const name = selOpt[i].name.split('-')[0];
      if (rcomp.includes(name)) continue;
      selOpt.splice(0, i);
      setOption((prevOption) => ({ ...prevOption, [name]: selOpt }));
      setRcomp((prevRcomp) => [...prevRcomp, name]);
      break;
    }
  }, [rcomp, selOpt]);

  const nextStep = () => {
    if (selected.length !== comp.length || currentProd === '') {
      console.log('Please select all components');
      msg.error('Please select all components');
      return false;
    }
    return true;
  }
  return (
    <Flex gap={5} align='center' justify='center' vertical>
      {contextHolder}
      <h1>Select a product</h1>
      {!prods ? (
        <p>No product available now</p>
      ) : (
        <Flex gap={10} align='center' justify='center' vertical>
          <Radio.Group defaultValue={null} onChange={ProdChange} size='large' buttonStyle='solid'>
            {prods.map((prod: string) => (
              <Radio.Button key={prod} value={prod} name={prod}>{prod}</Radio.Button>
            ))}</Radio.Group>
          <p>{prodDesc}</p>
          {rcomp.map((rcomp: string) => (
            <Flex>
              <Radio.Group key={rcomp} defaultValue={null} onChange={onSelect}>
              {option[rcomp] ? option[rcomp].map((value) => (
                <Popover content={value.desc} title={`${value.name} ($${value.price})`}>
                <Radio.Button key={value.name} value={value.name}>{value.name}</Radio.Button>
                </Popover>
              )): []}
              </Radio.Group>
            </Flex>
          ))}
          </Flex>
        )}
      <Result data={selected} render={nextStep}/>
    </Flex>
  );
};

export default App;