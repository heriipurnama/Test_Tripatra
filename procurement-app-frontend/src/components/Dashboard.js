import React, { useState } from 'react';
import { useQuery } from '@apollo/client';
import { gql } from '@apollo/client';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import CloudDownloadIcon from '@material-ui/icons/CloudDownload';
import TableContainer from '@material-ui/core/TableContainer';
import Table from '@material-ui/core/Table';
import TableHead from '@material-ui/core/TableHead';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import { jwtDecode } from "jwt-decode";
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import Drawer from '@material-ui/core/Drawer';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import HomeIcon from '@material-ui/icons/Home';
import ListAltIcon from '@material-ui/icons/ListAlt';
import { Link } from 'react-router-dom';
import LinearProgress from '@material-ui/core/LinearProgress';

const GET_PURCHASE_ORDERS = gql`
  query {
    getPurchaseOrders {
      orderId
      userId
      items {
        itemId
        name
        quantity
        price
      }
      totalAmount
      createdAt
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
  button: {
    margin: theme.spacing(2, 0),
  },
  tableContainer: {
    marginTop: theme.spacing(2),
  },
  logoutButton: {
    marginLeft: 'auto',
  },
  loadingBar: {
    width: '100%',
    '& > * + *': {
      marginTop: theme.spacing(2),
    },
  },
}));

const Dashboard = ({ onLogout }) => {
  const classes = useStyles();
  const { loading, error, data } = useQuery(GET_PURCHASE_ORDERS);
  const [user, setUser] = useState(getUserFromLocalStorage());

  function getUserFromLocalStorage() {
    const userData = localStorage.getItem('token');
    const tokenDecode = jwtDecode(userData);
    console.log("tokenDecode", tokenDecode);
    return tokenDecode ? tokenDecode : null;
  }

  const userDataLog  = getUserFromLocalStorage();
  console.log("userDataLog, ", userDataLog);
  
  const loadingBar = loading ? (
    <div className={classes.loadingBar}>
      <LinearProgress />
    </div>
  ) : null;

  if (error) return <p>Error: {error.message}</p>;

  const handleDownload = async () => {
    try {
      const response = await fetch('http://127.0.0.1:8080/download-purchase-orders');
      
      if (response.ok) {
        const responseData = await response.json();
        const base64PDF = responseData.base64PDF;
        const binaryPDF = atob(base64PDF);
        const arrayBufferPDF = new ArrayBuffer(binaryPDF.length);
        const uint8Array = new Uint8Array(arrayBufferPDF);
        for (let i = 0; i < binaryPDF.length; i++) {
          uint8Array[i] = binaryPDF.charCodeAt(i);
        }
        const blobPDF = new Blob([arrayBufferPDF], { type: 'application/pdf' });
        const url = URL.createObjectURL(blobPDF);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'purchase_orders.pdf';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
      } else {
        console.error('Error downloading PDF: Server response not OK');
      }
    } catch (error) {
      console.error('Error downloading PDF:', error);
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
          <div className={classes.logoutButton}>
            <IconButton color="inherit" onClick={handleDownload}>
              <CloudDownloadIcon />
            </IconButton>
            <IconButton color="inherit" onClick={handleLogout}>
              {user && (
                <Typography variant="body1" style={{ marginRight: '10px' }}>
                  {userDataLog.name}
                </Typography>
              )}
              <ExitToAppIcon />
            </IconButton>
          </div>
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
          <ListItem button >
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
          Dashboard
        </Typography>
        {loadingBar}
        <TableContainer component={Paper} className={classes.tableContainer}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Order ID</TableCell>
                <TableCell>User ID</TableCell>
                <TableCell>Items</TableCell>
                <TableCell>Total Amount</TableCell>
                <TableCell>Created At</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {data && data.getPurchaseOrders.map((order) => (
                <TableRow key={order.orderId}>
                  <TableCell>{order.orderId}</TableCell>
                  <TableCell>{order.userId}</TableCell>
                  <TableCell>
                    <ul>
                      {order.items.map((item) => (
                        <><li key={item.itemId}>Name: {item.name}</li>
                          <li key={item.itemId}>Quantity: {item.quantity}</li>
                          <li key={item.itemId}>Price: {item.price}</li>
                        </>
                      ))}
                    </ul>
                  </TableCell>
                  <TableCell>{order.totalAmount}</TableCell>
                  <TableCell>{order.createdAt}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </main>
    </div>
  );
};

export default Dashboard;
