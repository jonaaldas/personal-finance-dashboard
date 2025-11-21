export type PlaidAccountBalance = {
  available: number | null;
  current: number;
  iso_currency_code: string | null;
  limit: number | null;
  unofficial_currency_code: string | null;
};

export type PlaidAccount = {
  account_id: string;
  balances: PlaidAccountBalance;
  holder_category?: string | null;
  mask: string;
  name: string;
  official_name: string | null;
  subtype: string;
  type: string;
};

export type PlaidItem = {
  auth_method?: string | null;
  available_products?: string[];
  billed_products?: string[];
  consent_expiration_time?: string | null;
  consented_products?: string[];
  error?: any | null;
  institution_id: string;
  institution_name: string;
  institution_logo?: string | null;
  institution_url?: string | null;
  item_id: string;
  products?: string[];
  update_type?: string;
  webhook?: string;
};

export type PlaidAccountsGetResponse = {
  accounts: PlaidAccount[];
  item: PlaidItem;
  request_id: string;
};

export type SaveAccountsResponse = {
  message: string;
};

