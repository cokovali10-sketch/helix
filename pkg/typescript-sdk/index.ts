// Helix TypeScript SDK
// This SDK allows TypeScript/JavaScript applications to interact with the Helix observability platform

/**
 * Client - Main Helix client for interacting with the platform
 */
export class Client {
  private baseURL: string;
  private apiKey: string;
  private httpClient: XMLHttpRequest | fetch;

  /**
   * Constructor
   * @param baseURL - The Helix API base URL (e.g., "http://localhost:8080")
   * @param apiKey - The API authentication key
   */
  constructor(baseURL: string, apiKey: string) {
    this.baseURL = baseURL;
    this.apiKey = apiKey;
    this.httpClient = fetch; // Using fetch API
  }

  /**
   * Ingest an event into the Helix platform
   * @param event - The event to ingest
   * @returns Promise resolving to the ingestion response
   */
  async ingestEvent(event: Event): Promise<IngestionResponse> {
    const response = await fetch(`${this.baseURL}/v1/events`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.apiKey}`,
      },
      body: JSON.stringify(event),
    });

    if (!response.ok) {
      throw new Error(`Server returned status ${response.status}`);
    }

    return response.json();
  }

  /**
   * Execute a query against the Helix platform
   * @param options - Query options including expression and tenant
   * @returns Promise resolving to query results
   */
  async query(options: QueryOptions): Promise<QueryResult[]> {
    const response = await fetch(`${this.baseURL}/v1/query`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.apiKey}`,
      },
      body: JSON.stringify(options),
    });

    if (!response.ok) {
      throw new Error(`Server returned status ${response.status}`);
    }

    return response.json();
  }

  /**
   * List all tenants
   * @returns Promise resolving to tenant information
   */
  async listTenants(): Promise<TenantInfo[]> {
    const response = await fetch(`${this.baseURL}/v1/tenants`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
      },
    });

    if (!response.ok) {
      throw new Error(`Server returned status ${response.status}`);
    }

    return response.json();
  }

  /**
   * Check the health of the Helix platform
   * @returns Promise resolving to boolean indicating health status
   */
  async healthCheck(): Promise<boolean> {
    const response = await fetch(`${this.baseURL}/health`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
      },
    });

    return response.ok;
  }

  /**
   * Get platform metrics
   * @returns Promise resolving to metrics object
   */
  async metrics(): Promise<Metrics> {
    const response = await fetch(`${this.baseURL}/metrics`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
      },
    });

    if (!response.ok) {
      throw new Error(`Server returned status ${response.status}`);
    }

    return response.json();
  }

  /**
   * Search for events with optional filters
   * @param query - The search query string
   * @param filters - Optional key-value filter pairs
   * @returns Promise resolving to matching events
   */
  async searchEvents(query: string, filters?: Record<string, string>): Promise<Event[]> {
    let url = `${this.baseURL}/v1/events/search?query=${encodeURIComponent(query)}`;

    if (filters) {
      const filterParams = Object.entries(filters)
        .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
        .join('&');
      url += `&${filterParams}`;
    }

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
      },
    });

    if (!response.ok) {
      throw new Error(`Server returned status ${response.status}`);
    }

    return response.json();
  }
}

/**
 * Event - An observability event
 */
export interface Event {
  id?: string;
  timestamp: number; // Unix timestamp in milliseconds
  source: string;
  metrics: Record<string, number>;
  tags: Record<string, string>;
  tenant: string;
}

/**
 * IngestionResponse - Response from event ingestion
 */
export interface IngestionResponse {
  success: boolean;
  message: string;
  eventId: string;
}

/**
 * QueryOptions - Options for executing a query
 */
export interface QueryOptions {
  expression: string;
  tenant?: string;
  timeRange?: [number, number]; // [start, end] Unix timestamps in milliseconds
}

/**
 * QueryResult - Result of a single query
 */
export interface QueryResult {
  metric: Record<string, string>;
  value: number;
  timestamp: number;
}

/**
 * TenantInfo - Tenant information
 */
export interface TenantInfo {
  id: string;
  name: string;
  rateLimit: number;
  quota: number;
}

/**
 * Metrics - Platform metrics
 */
export interface Metrics {
  eventsReceived: number;
  eventsProcessed: number;
  activeTenants: number;
  memoryUsageBytes: number;
}

/**
 * Helper function to create a client instance
 * @param baseURL - The Helix API base URL
 * @param apiKey - The API authentication key
 * @returns A new Client instance
 */
export function createClient(baseURL: string, apiKey: string): Client {
  return new Client(baseURL, apiKey);
}

/**
 * Pre-configured client for local development
 */
export const defaultClient = createClient('http://localhost:8080', 'helix-dev-key');