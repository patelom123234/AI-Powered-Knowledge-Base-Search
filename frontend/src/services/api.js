const API_BASE_URL = 'http://localhost:8080';

export const postSearchQuery = async (query) => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/search-query`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ query: query }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    return data;

  } catch (error) {
    console.error("Error in postSearchQuery:", error);
    throw error;
  }
};