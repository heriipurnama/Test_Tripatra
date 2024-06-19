// PurchaseOrderForm.jsx
import React, { useState } from 'react';
import { useMutation } from '@apollo/client';
import { gql } from '@apollo/client';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import { jwtDecode } from "jwt-decode";
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import HomeIcon from '@material-ui/icons/Home';
import ListAltIcon from '@material-ui/icons/ListAlt';
import { Link } from 'react-router-dom';

const CREATE_PURCHASE_ORDER = gql`
  mutation CreatePurchaseOrder($userId: ID!, $items: [ItemInput!]!) {
    createPurchaseOrder(userId: $userId, items: $items) {
      orderId
      userId
      items {
        itemId
        name
        quantity
        price
      }
      totalAmount
    }
  }
`;

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing(3),
  },
  toolbar: theme.mixins.toolbar,
  form: {
    '& > div': {
      marginBottom: theme.spacing(2),
    },
  },
  logoutButton: {
    marginLeft: 'auto',
  },
}));

const PurchaseOrderForm = ({ onLogout }) => {
  const classes = useStyles();
  const [items, setItems] = useState([]);
  const [createPurchaseOrder] = useMutation(CREATE_PURCHASE_ORDER);
  const [user, setUser] = useState(getUserFromLocalStorage());

  const handleAddItem = () => {
    setItems([...items, { itemId: '', name: '', quantity: 0, price: 0 }]);
  };

  function getUserFromLocalStorage() {
    const userData = localStorage.getItem('token');
    const tokenDecode = jwtDecode(userData);
    console.log("tokenDecode", tokenDecode);
    return tokenDecode ? tokenDecode : null;
  }

  const userDataLog = getUserFromLocalStorage();
  console.log("userDataLog, ", userDataLog);

  const handleChangeItem = (index, field, value) => {
    const updatedItems = [...items];
    updatedItems[index][field] = value;
    setItems(updatedItems);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const { data } = await createPurchaseOrder({ variables: { userId: userDataLog.name, items } });
      console.log('Purchase Order Created:', data.createPurchaseOrder);
    } catch (error) {
      console.error('Error creating purchase order:', error);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('user');
    setUser(null);
    onLogout();
  };

  return (
    <div className={classes.root}>
      <AppBar position="fixed" className={classes.appBar}>
        <Toolbar>
          <Typography variant="h6" noWrap>
            Procurement App
          </Typography>
          <IconButton color="inherit" className={classes.logoutButton} onClick={handleLogout}>
            {user && (
              <Typography variant="body1" style={{ marginRight: '10px' }}>
                {userDataLog.name}
              </Typography>
            )}
            <ExitToAppIcon />
          </IconButton>
        </Toolbar>
      </AppBar>
      <Drawer
        className={classes.drawer}  
        variant="permanent"
        classes={{
          paper: classes.drawerPaper,
        }}
      >
        <div className={classes.toolbar} />
        <Divider />
        <List>
          <ListItem button>
            <ListItemIcon><HomeIcon /></ListItemIcon>
            <ListItemText primary="Home" />
          </ListItem>
          <ListItem button component={Link} to="/dashboard">
            <ListItemIcon><ListAltIcon /></ListItemIcon>
            <ListItemText primary="Dashboard" />
          </ListItem>
          <ListItem button component={Link} to="/purchase-order-form">
            <ListItemIcon><ListAltIcon /></ListItemIcon>
            <ListItemText primary="Purchase Orders" />
          </ListItem>
        </List>
        <Divider />
      </Drawer>
      <main className={classes.content}>
        <div className={classes.toolbar} />
        <Typography variant="h2" component="h2" gutterBottom>
          Create Purchase Order
        </Typography>
        <form className={classes.form} onSubmit={handleSubmit}>
          <TextField
            label="User ID"
            variant="outlined"
            value={userDataLog.name}
            readOnly={true}
          />
          {items.map((item, index) => (
            <div key={index}>
              <TextField
                label="Item ID"
                variant="outlined"
                value={item.itemId}
                onChange={(e) => handleChangeItem(index, 'itemId', e.target.value)}
              />
              <TextField
                label="Name"
                variant="outlined"
                value={item.name}
                onChange={(e) => handleChangeItem(index, 'name', e.target.value)}
              />
              <TextField
                label="Quantity"
                variant="outlined"
                type="number"
                value={item.quantity}
                onChange={(e) => handleChangeItem(index, 'quantity', parseInt(e.target.value))}
              />
              <TextField
                label="Price"
                variant="outlined"
                type="number"
                value={item.price}
                onChange={(e) => handleChangeItem(index, 'price', parseFloat(e.target.value))}
              />
            </div>
          ))}
          <Button variant="contained" color="primary" onClick={handleAddItem}>Add Item</Button>
          <Button variant="contained" color="primary" type="submit">Submit</Button>
        </form>
      </main>
    </div>
  );
};

export default PurchaseOrderForm;
