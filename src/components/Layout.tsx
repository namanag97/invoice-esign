import React, { useState } from 'react';
import { Sidebar, SidebarProps } from './Sidebar';
import { Header, HeaderProps } from './Header';
import './Layout.css';

export interface LayoutProps {
  /**
   * Props for the Sidebar component
   */
  sidebarProps: Omit<SidebarProps, 'collapsed' | 'onToggleCollapse'>;
  /**
   * Props for the Header component
   */
  headerProps: HeaderProps;
  /**
   * Child content to render in the main area
   */
  children: React.ReactNode;
  /**
   * Initial sidebar collapsed state
   */
  initialSidebarCollapsed?: boolean;
  /**
   * Additional CSS class for the layout
   */
  className?: string;
}

/**
 * Main application layout with sidebar and header
 */
export const Layout = ({
  sidebarProps,
  headerProps,
  children,
  initialSidebarCollapsed = false,
  className = '',
}: LayoutProps) => {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(initialSidebarCollapsed);

  const toggleSidebar = () => {
    setSidebarCollapsed(prev => !prev);
  };

  return (
    <div className={`layout ${className}`}>
      <Sidebar
        {...sidebarProps}
        collapsed={sidebarCollapsed}
        onToggleCollapse={toggleSidebar}
      />
      <div className={`layout-content ${sidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <Header {...headerProps} />
        <main className="layout-main">
          {children}
        </main>
      </div>
    </div>
  );
}; 