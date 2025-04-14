import React, { InputHTMLAttributes } from 'react';
import './Input.css';

export interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  /**
   * Label for the input
   */
  label?: string;
  /**
   * Input size
   */
  size?: 'small' | 'medium' | 'large';
  /**
   * Input variant
   */
  variant?: 'outlined' | 'filled' | 'standard';
  /**
   * Error message
   */
  error?: string;
  /**
   * Helper text
   */
  helperText?: string;
}

/**
 * Input UI component for user text input
 */
export const Input = ({
  label,
  size = 'medium',
  variant = 'outlined',
  error,
  helperText,
  className = '',
  ...props
}: InputProps) => {
  const inputSize = `input--${size}`;
  const inputVariant = `input--${variant}`;
  const inputError = error ? 'input--error' : '';
  
  return (
    <div className="input-container">
      {label && <label className="input-label">{label}</label>}
      <input
        className={`input ${inputSize} ${inputVariant} ${inputError} ${className}`}
        {...props}
      />
      {helperText && <p className="input-helper-text">{helperText}</p>}
      {error && <p className="input-error-text">{error}</p>}
    </div>
  );
}; 