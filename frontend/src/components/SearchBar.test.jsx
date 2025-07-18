import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';

import SearchBar from './SearchBar';

describe('SearchBar', () => {

  it('should render an input field and a submit button', () => {
    render(<SearchBar onSearch={() => {}} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask anything about our it systems/i);
    const buttonElement = screen.getByRole('button', { name: /search/i });

    expect(inputElement).toBeInTheDocument();
    expect(buttonElement).toBeInTheDocument();
  });

  it('should render search suggestions', () => {
    render(<SearchBar onSearch={() => {}} isLoading={false} />);

    const suggestionsText = screen.getByText('Try asking about:');
    const passwordResetTag = screen.getByRole('button', { name: /password reset/i });
    const vpnSetupTag = screen.getByRole('button', { name: /vpn setup/i });
    const printerSetupTag = screen.getByRole('button', { name: /printer setup/i });

    expect(suggestionsText).toBeInTheDocument();
    expect(passwordResetTag).toBeInTheDocument();
    expect(vpnSetupTag).toBeInTheDocument();
    expect(printerSetupTag).toBeInTheDocument();
  });

  it('should call the onSearch prop with the query when the form is submitted', () => {
    const mockOnSearch = vi.fn();
    render(<SearchBar onSearch={mockOnSearch} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask anything about our it systems/i);
    const buttonElement = screen.getByRole('button', { name: /search/i });

    fireEvent.change(inputElement, { target: { value: 'test query' } });
    fireEvent.click(buttonElement);

    expect(mockOnSearch).toHaveBeenCalledTimes(1);
    expect(mockOnSearch).toHaveBeenCalledWith('test query');
  });

  it('should call onSearch when Enter key is pressed', () => {
    const mockOnSearch = vi.fn();
    render(<SearchBar onSearch={mockOnSearch} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask anything about our it systems/i);
    
    fireEvent.change(inputElement, { target: { value: 'test query' } });
    fireEvent.keyDown(inputElement, { key: 'Enter', code: 'Enter' });

    expect(mockOnSearch).toHaveBeenCalledTimes(1);
    expect(mockOnSearch).toHaveBeenCalledWith('test query');
  });

  it('should not call onSearch when Enter key is pressed with Shift key', () => {
    const mockOnSearch = vi.fn();
    render(<SearchBar onSearch={mockOnSearch} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask anything about our it systems/i);
    
    fireEvent.change(inputElement, { target: { value: 'test query' } });
    fireEvent.keyDown(inputElement, { key: 'Enter', code: 'Enter', shiftKey: true });

    expect(mockOnSearch).not.toHaveBeenCalled();
  });

  it('should not call onSearch when form is submitted with empty query', () => {
    const mockOnSearch = vi.fn();
    render(<SearchBar onSearch={mockOnSearch} isLoading={false} />);

    const buttonElement = screen.getByRole('button', { name: /search/i });
    fireEvent.click(buttonElement);

    expect(mockOnSearch).not.toHaveBeenCalled();
  });

  it('should not call onSearch when Enter key is pressed with empty query', () => {
    const mockOnSearch = vi.fn();
    render(<SearchBar onSearch={mockOnSearch} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask anything about our it systems/i);
    fireEvent.keyDown(inputElement, { key: 'Enter', code: 'Enter' });

    expect(mockOnSearch).not.toHaveBeenCalled();
  });

  it('should disable the input and button when isLoading is true', () => {
    render(<SearchBar onSearch={() => {}} isLoading={true} />);

    const inputElement = screen.getByPlaceholderText(/ask anything about our it systems/i);
    const buttonElement = screen.getByRole('button', { name: /searching.../i });

    expect(inputElement).toBeDisabled();
    expect(buttonElement).toBeDisabled();
  });

  it('should disable suggestion tags when isLoading is true', () => {
    render(<SearchBar onSearch={() => {}} isLoading={true} />);

    const passwordResetTag = screen.getByRole('button', { name: /password reset/i });
    const vpnSetupTag = screen.getByRole('button', { name: /vpn setup/i });
    const printerSetupTag = screen.getByRole('button', { name: /printer setup/i });

    expect(passwordResetTag).toBeDisabled();
    expect(vpnSetupTag).toBeDisabled();
    expect(printerSetupTag).toBeDisabled();
  });

  it('should disable submit button when query is empty', () => {
    render(<SearchBar onSearch={() => {}} isLoading={false} />);

    const buttonElement = screen.getByRole('button', { name: /search/i });
    expect(buttonElement).toBeDisabled();
  });

  it('should enable submit button when query has content', () => {
    render(<SearchBar onSearch={() => {}} isLoading={false} />);

    const inputElement = screen.getByPlaceholderText(/ask anything about our it systems/i);
    const buttonElement = screen.getByRole('button', { name: /search/i });

    fireEvent.change(inputElement, { target: { value: 'test query' } });
    expect(buttonElement).not.toBeDisabled();
  });

  it('should call onSearch with specific query when suggestion tags are clicked', () => {
    const mockOnSearch = vi.fn();
    render(<SearchBar onSearch={mockOnSearch} isLoading={false} />);

    const passwordResetTag = screen.getByRole('button', { name: /password reset/i });
    const vpnSetupTag = screen.getByRole('button', { name: /vpn setup/i });
    const printerSetupTag = screen.getByRole('button', { name: /printer setup/i });

    fireEvent.click(passwordResetTag);
    expect(mockOnSearch).toHaveBeenCalledWith('How do I reset my password?');

    fireEvent.click(vpnSetupTag);
    expect(mockOnSearch).toHaveBeenCalledWith('VPN connection issues');

    fireEvent.click(printerSetupTag);
    expect(mockOnSearch).toHaveBeenCalledWith('Setting up a new printer');
  });

  it('should show loading state in button when isLoading is true', () => {
    render(<SearchBar onSearch={() => {}} isLoading={true} />);

    const loadingButton = screen.getByRole('button', { name: /searching.../i });
    expect(loadingButton).toBeInTheDocument();
    expect(loadingButton).toHaveClass('loading');
  });

  it('should show search icon and arrow in button when not loading', () => {
    render(<SearchBar onSearch={() => {}} isLoading={false} />);

    const searchButton = screen.getByRole('button', { name: /search/i });
    expect(searchButton).toBeInTheDocument();
    expect(searchButton).not.toHaveClass('loading');
  });

  it('should render search icon in input container', () => {
    render(<SearchBar onSearch={() => {}} isLoading={false} />);

    const searchIcon = screen.getByText('🔍');
    expect(searchIcon).toBeInTheDocument();
  });
});