import React from 'react';
import { Link, useLocation } from 'react-router-dom'; // Import useLocation
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
}));

const LeftBar = () => {
  const classes = useStyles();
  const location = useLocation(); // Get the current location

  // Determine whether to hide the LeftBar based on the current route
  const isSignInPage = location.pathname === '/signin';

  if (isSignInPage) {
    return null; // If it's the SignIn page, don't render the LeftBar
  }

    // Determine whether to hide the LeftBar based on the current route
    const isSignUpPage = location.pathname === '/signup';

    if (isSignUpPage) {
      return null; // If it's the SignIn page, don't render the LeftBar
    }

  return (
    <Drawer
      className={classes.drawer}
      variant="permanent"
      classes={{
        paper: classes.drawerPaper,
      }}
    >
      <List>
        <ListItem button component={Link} to="/dashboard">
          <ListItemText primary="Dashboard" />
        </ListItem>
        <ListItem button component={Link} to="/purchase-order-form">
          <ListItemText primary="Purchase Order Form" />
        </ListItem>
        {/* Add more menu items here */}
      </List>
    </Drawer>
  );
};

export default LeftBar;
