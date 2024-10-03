import { Space } from "antd";
import { useSelector, useDispatch } from "react-redux";
import { increment, decrement, selectCount } from './counterSlice'
import { MinusCircleOutlined, PlusCircleOutlined } from "@ant-design/icons";

export default function Counter(){
    const count = useSelector(selectCount);
    const dispatch = useDispatch();
    return(
        <Space direction="horizontal">
        <a onClick={() => dispatch(decrement(count))} target="_blank"><MinusCircleOutlined /></a>
        {count}
        <a onClick={() => dispatch(increment(count))} target="_blank"><PlusCircleOutlined /></a>
        </Space>
    )
}