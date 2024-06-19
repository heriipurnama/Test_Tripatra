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
import { onError } from '@apollo/client/link/error'; 
import LeftBar from './components/LeftBar'; // import LeftBar component

const httpLink = createHttpLink({
  uri: 'http://localhost:8080/graphql',
});

const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : "",
    }
  }
});

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path }) =>
      console.log(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`)
    );
  }

  if (networkError) {
    console.log(`[Network error]: ${networkError}`);
  }
});

// Gabungkan link dengan urutan: errorLink -> authLink -> httpLink
const link = errorLink.concat(authLink).concat(httpLink);

const client = new ApolloClient({
  link: link,  // Gunakan link gabungan
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
        <LeftBar />
          <Routes>
            <Route path="/signin" element={<SignIn onLogin={handleLogin} />} />
            <Route path="/signup" element={<SignUp />} />
            <Route
              path="/dashboard"
              element={loggedIn ? <Dashboard onLogout={handleLogout} /> : <Navigate to="/signin" />}
            />
            <Route
              path="/purchase-order-form"
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
