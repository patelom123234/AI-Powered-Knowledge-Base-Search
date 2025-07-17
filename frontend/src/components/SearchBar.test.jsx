import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';

import SearchBar from './SearchBar';

describe('SearchBar', () => {

  it('should render an input field and a submit button', () => {
    render(<SearchBar onSearch={() => {}} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask a question/i);
    const buttonElement = screen.getByRole('button', { name: /search/i });

    expect(inputElement).toBeInTheDocument();
    expect(buttonElement).toBeInTheDocument();
  });


  it('should call the onSearch prop with the query when the form is submitted', () => {
    const mockOnSearch = vi.fn();
    render(<SearchBar onSearch={mockOnSearch} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask a question/i);
    const buttonElement = screen.getByRole('button', { name: /search/i });

    fireEvent.change(inputElement, { target: { value: 'test query' } });

    fireEvent.click(buttonElement);

    expect(mockOnSearch).toHaveBeenCalledTimes(1);
    expect(mockOnSearch).toHaveBeenCalledWith('test query');
  });

  it('should disable the input and button when isLoading is true', () => {
    render(<SearchBar onSearch={() => {}} isLoading={true} />);

    const inputElement = screen.getByPlaceholderText(/ask a question/i);
    const buttonElement = screen.getByRole('button', { name: /searching.../i });

    expect(inputElement).toBeDisabled();
    expect(buttonElement).toBeDisabled();
  });
});