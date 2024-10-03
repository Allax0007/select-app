import { selectToken, logout } from './AuthSlice';
import { useSelector, useDispatch } from 'react-redux';

export const authenticate = async () => {
  const dispatch = useDispatch();
  const token = useSelector(selectToken);
  try {  
    const response = await fetch('/validate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      throw new Error('Failed to authenticate');
    }
  } catch (error: any) {
    console.error('Error authenticating:', error);
    dispatch(logout(token));
  }
};