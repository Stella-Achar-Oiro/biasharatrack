import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './homepage/Home';
import Login from './pages/Login';
import SignUp from './pages/SignUp';
import NotFound from './pages/error/NotFound';
import Unauthorized from './pages/error/Unauthorized';
import ServerError from './pages/error/ServerError';
import ErrorBoundary from './homepage/components/ErrorBoundary';
import FloatingLanguageSwitcher from './homepage/components/FloatingLanguageSwitcher';
import React from 'react';


function App() {
  return (
    <ErrorBoundary>
      <Router>
        <FloatingLanguageSwitcher />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<SignUp />} />
          
          <Route path="/unauthorized" element={<Unauthorized />} />
          <Route path="/server-error" element={<ServerError />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Router>
    </ErrorBoundary>
  );
}

export default App;