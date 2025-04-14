import React, { useState } from 'react';
import './Sidebar.css';

export interface SidebarItem {
  /**
   * Item label
   */
  label: string;
  /**
   * Icon component or element
   */
  icon?: React.ReactNode;
  /**
   * URL or callback function for item
   */
  action?: string | (() => void);
  /**
   * Whether the item is active
   */
  active?: boolean;
  /**
   * Optional badge or notification count
   */
  badge?: number | string;
  /**
   * Nested children items
   */
  children?: SidebarItem[];
}

export interface SidebarProps {
  /**
   * Company/product logo
   */
  logo: React.ReactNode;
  /**
   * Sidebar navigation items
   */
  items: SidebarItem[];
  /**
   * Whether sidebar is collapsed
   */
  collapsed?: boolean;
  /**
   * Toggle collapsed state callback
   */
  onToggleCollapse?: () => void;
  /**
   * User profile element at bottom of sidebar
   */
  userProfile?: React.ReactNode;
  /**
   * Additional CSS class
   */
  className?: string;
}

/**
 * YC SaaS style sidebar for application navigation
 */
export const Sidebar = ({
  logo,
  items,
  collapsed = false,
  onToggleCollapse,
  userProfile,
  className = '',
}: SidebarProps) => {
  const [expandedItems, setExpandedItems] = useState<Record<string, boolean>>({});

  const toggleExpand = (label: string) => {
    setExpandedItems(prev => ({
      ...prev,
      [label]: !prev[label]
    }));
  };

  const renderItem = (item: SidebarItem, index: number) => {
    const hasChildren = item.children && item.children.length > 0;
    const isExpanded = expandedItems[item.label];
    
    return (
      <li key={index} className={`sidebar-item ${item.active ? 'active' : ''}`}>
        <div 
          className="sidebar-item-content"
          onClick={() => {
            if (hasChildren) {
              toggleExpand(item.label);
            } else if (typeof item.action === 'function') {
              item.action();
            }
          }}
        >
          {item.icon && <span className="sidebar-item-icon">{item.icon}</span>}
          {!collapsed && (
            <>
              <span className="sidebar-item-label">{item.label}</span>
              {item.badge && <span className="sidebar-item-badge">{item.badge}</span>}
              {hasChildren && (
                <span className={`sidebar-item-arrow ${isExpanded ? 'expanded' : ''}`}>
                  ▾
                </span>
              )}
            </>
          )}
        </div>
        
        {hasChildren && isExpanded && !collapsed && (
          <ul className="sidebar-submenu">
            {item.children!.map((child, childIndex) => (
              <li 
                key={childIndex} 
                className={`sidebar-submenu-item ${child.active ? 'active' : ''}`}
                onClick={() => {
                  if (typeof child.action === 'function') {
                    child.action();
                  }
                }}
              >
                {child.icon && <span className="sidebar-item-icon">{child.icon}</span>}
                <span className="sidebar-item-label">{child.label}</span>
                {child.badge && <span className="sidebar-item-badge">{child.badge}</span>}
              </li>
            ))}
          </ul>
        )}
      </li>
    );
  };

  return (
    <div className={`sidebar ${collapsed ? 'collapsed' : ''} ${className}`}>
      <div className="sidebar-header">
        <div className="sidebar-logo">
          {logo}
        </div>
        {onToggleCollapse && (
          <button 
            className="sidebar-toggle"
            onClick={onToggleCollapse}
            aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
          >
            {collapsed ? '→' : '←'}
          </button>
        )}
      </div>
      
      <nav className="sidebar-nav">
        <ul className="sidebar-menu">
          {items.map(renderItem)}
        </ul>
      </nav>
      
      {userProfile && (
        <div className="sidebar-footer">
          {userProfile}
        </div>
      )}
    </div>
  );
}; 