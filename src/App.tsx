import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';

import AdminLayout from './layouts/AdminLayout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Players from './pages/Players';
import Admins from './pages/Admins';
import FishConfig from './pages/FishConfig';
import Rooms from './pages/Rooms';
import Settings from './pages/Settings';
import GamePlay from './pages/GamePlay';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        
        {/* Protected Routes enclosed in AdminLayout */}
        <Route path="/" element={<AdminLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="players" element={<Players />} />
          <Route path="admins" element={<Admins />} />
          <Route path="fish" element={<FishConfig />} />
          <Route path="rooms" element={<Rooms />} />
          <Route path="settings" element={<Settings />} />
          <Route path="gameplay" element={<GamePlay />} />
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
