import { useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { logout, selectExp } from '../features/Auth/AuthSlice';
import { Button, Modal } from 'antd';

function Timer() {
  const [open, setOpen] = useState(false);
  const [timer, setTimer] = useState(0);
  const dispatch = useDispatch();
  var exp = useSelector(selectExp);

  useEffect(() => {
    const interval = setInterval(() => {
      if (exp && exp < Math.floor(Date.now() / 1000)) {
        dispatch(logout(exp));
        setOpen(true);
      }
      setTimer(exp - Math.floor(Date.now() / 1000));
    }, 0);
    return () => clearInterval(interval);
  }, [exp]);

  return (
    <>
    <>{exp ? timer : exp}</>
    <Modal 
      title='Session expired'
      open={open}
      closable={false}
      footer={<Button type='primary' onClick={()=>setOpen(false)}>Ok</Button>}
      >
      <p>Please log in again.</p>
    </Modal>
    </>
  );
}

export default Timer;
