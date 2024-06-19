import React, { useState } from 'react';
import { ApolloProvider, ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import { Container, CssBaseline } from '@material-ui/core';
import SignIn from './components/SignIn';
import SignUp from './components/SignUp';
import Dashboard from './components/Dashboard'; // Import Dashboard component
import PurchaseOrderForm from './components/PurchaseOrderForm'; // Import PurchaseOrderForm component
import ReportGenerator from './components/ReportGenerator'; // Import ReportGenerator component

// Apollo GraphQL client setup
const httpLink = createHttpLink({
  uri: 'http://127.0.0.1:8080/grapql', // GraphQL server URL
});

const authLink = setContext((_, { headers }) => {
  // Get token from localStorage
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '', // Attach token to request headers
    }
  }
});

const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache()
});

const App = () => {
  const [loggedIn, setLoggedIn] = useState(!!localStorage.getItem('token'));

  const handleLogin = () => {
    setLoggedIn(true);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setLoggedIn(false);
  };

  return (
    <ApolloProvider client={client}>
      <Router>
        <CssBaseline />
        <Container>
          <Routes>
            <Route path="/signin" element={<SignIn onLogin={handleLogin} />} />
            <Route path="/signup" element={<SignUp />} />
            <Route
              path="/dashboard"
              element={loggedIn ? <Dashboard onLogout={handleLogout} /> : <Navigate to="/signin" />}
            />
            <Route
              path="/create-order"
              element={loggedIn ? <PurchaseOrderForm /> : <Navigate to="/signin" />}
            />
            <Route
              path="/generate-report"
              element={loggedIn ? <ReportGenerator /> : <Navigate to="/signin" />}
            />
            <Route path="/" element={<Navigate to="/signin" />} />
          </Routes>
        </Container>
      </Router>
    </ApolloProvider>
  );
};

export default App;
