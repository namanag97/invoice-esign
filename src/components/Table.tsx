import React from 'react';
import './Table.css';

export interface TableColumn<T> {
  header: string;
  accessor: keyof T | ((data: T) => React.ReactNode);
  className?: string;
}

export interface TableProps<T> {
  /**
   * Table columns definition
   */
  columns: TableColumn<T>[];
  /**
   * Data to be displayed in the table
   */
  data: T[];
  /**
   * Optional CSS class for the table
   */
  className?: string;
  /**
   * Optional row click handler
   */
  onRowClick?: (item: T, index: number) => void;
  /**
   * Whether to show a hover effect on rows
   */
  hoverEffect?: boolean;
  /**
   * Whether table has fixed header
   */
  fixedHeader?: boolean;
  /**
   * Maximum height for the table with fixed header
   */
  maxHeight?: string;
}

/**
 * Modern table component for displaying structured data
 */
export const Table = <T extends Record<string, any>>({
  columns,
  data,
  className = '',
  onRowClick,
  hoverEffect = true,
  fixedHeader = false,
  maxHeight = '600px',
}: TableProps<T>) => {
  return (
    <div 
      className={`table-container ${fixedHeader ? 'table-fixed-header' : ''}`}
      style={fixedHeader ? { maxHeight } : undefined}
    >
      <table className={`table ${hoverEffect ? 'table-hover' : ''} ${className}`}>
        <thead>
          <tr>
            {columns.map((column, index) => (
              <th key={index} className={column.className}>
                {column.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length > 0 ? (
            data.map((item, rowIndex) => (
              <tr 
                key={rowIndex}
                onClick={onRowClick ? () => onRowClick(item, rowIndex) : undefined}
                className={onRowClick ? 'clickable' : ''}
              >
                {columns.map((column, colIndex) => {
                  const cellValue = typeof column.accessor === 'function'
                    ? column.accessor(item)
                    : item[column.accessor];
                  
                  return (
                    <td key={colIndex} className={column.className}>
                      {cellValue}
                    </td>
                  );
                })}
              </tr>
            ))
          ) : (
            <tr className="table-empty">
              <td colSpan={columns.length}>No data available</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}; 