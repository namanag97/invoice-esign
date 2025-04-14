import React, { ReactNode } from 'react';
import './Card.css';

export interface CardProps {
  /**
   * Card title
   */
  title: string;
  /**
   * Card content
   */
  children: ReactNode;
  /**
   * Card width
   */
  width?: string;
  /**
   * Optional border color
   */
  borderColor?: string;
}

/**
 * Card UI component for displaying content in a contained area
 */
export const Card = ({
  title,
  children,
  width = '300px',
  borderColor,
  ...props
}: CardProps) => {
  return (
    <div 
      className="card"
      style={{ 
        width, 
        borderColor: borderColor || undefined 
      }}
      {...props}
    >
      <div className="card-header">
        <h3>{title}</h3>
      </div>
      <div className="card-content">
        {children}
      </div>
    </div>
  );
}; 