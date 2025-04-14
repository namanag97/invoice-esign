import React from 'react';
import './Header.css';

export interface HeaderAction {
  /**
   * Label for the action
   */
  label?: string;
  /**
   * Icon element for the action
   */
  icon?: React.ReactNode;
  /**
   * Callback function for the action
   */
  onClick?: () => void;
  /**
   * Additional CSS class for the action
   */
  className?: string;
}

export interface HeaderProps {
  /**
   * Page title
   */
  title?: string;
  /**
   * Subtitle or breadcrumb
   */
  subtitle?: React.ReactNode;
  /**
   * Actions displayed on the right side
   */
  actions?: HeaderAction[];
  /**
   * Search component
   */
  searchComponent?: React.ReactNode;
  /**
   * Notification component (icon or button)
   */
  notificationComponent?: React.ReactNode;
  /**
   * User profile/avatar component
   */
  userComponent?: React.ReactNode;
  /**
   * Additional CSS class for header
   */
  className?: string;
  /**
   * Whether the header is sticky
   */
  sticky?: boolean;
  /**
   * Whether to show a border
   */
  bordered?: boolean;
}

/**
 * Modern SaaS app header component
 */
export const Header = ({
  title,
  subtitle,
  actions = [],
  searchComponent,
  notificationComponent,
  userComponent,
  className = '',
  sticky = true,
  bordered = true,
}: HeaderProps) => {
  return (
    <header 
      className={`header ${sticky ? 'sticky' : ''} ${bordered ? 'bordered' : ''} ${className}`}
    >
      <div className="header-left">
        {(title || subtitle) && (
          <div className="header-title-container">
            {title && <h1 className="header-title">{title}</h1>}
            {subtitle && <div className="header-subtitle">{subtitle}</div>}
          </div>
        )}
      </div>
      
      <div className="header-right">
        {actions.length > 0 && (
          <div className="header-actions">
            {actions.map((action, index) => (
              <button
                key={index}
                className={`header-action ${action.className || ''}`}
                onClick={action.onClick}
                aria-label={action.label}
              >
                {action.icon}
                {action.label && <span>{action.label}</span>}
              </button>
            ))}
          </div>
        )}
        
        {searchComponent && (
          <div className="header-search">
            {searchComponent}
          </div>
        )}
        
        {notificationComponent && (
          <div className="header-notification">
            {notificationComponent}
          </div>
        )}
        
        {userComponent && (
          <div className="header-user">
            {userComponent}
          </div>
        )}
      </div>
    </header>
  );
}; 