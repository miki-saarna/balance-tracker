export interface FilterOptions {
  accountType: AccountTypes;
  balance: AccountBalance;
}

export type AccountTypesKeys = Array<keyof AccountTypes>; // (keyof AccountTypes)[];

export type AccountTypes = Partial<Record<PlaidAccountTypes, boolean>>;
// PlaidAccountTypes?: boolean;
// [key in PlaidAccountTypes]: boolean;

type AccountBalance = {
  equalTo?: number;
  greaterThan?: number;
  lessThan?: number;
};

export type PlaidAccountTypes =
  | "401a"
  | "401k"
  | "403B"
  | "457b"
  | "529"
  | "brokerage"
  | "cash isa"
  | "crypto exchange"
  | "education savings account"
  | "ebt"
  | "fixed annuity"
  | "gic"
  | "health reimbursement arrangement"
  | "hsa"
  | "isa"
  | "ira"
  | "lif"
  | "life insurance"
  | "lira"
  | "lrif"
  | "lrsp"
  | "non-custodial wallet"
  | "non-taxable brokerage account"
  | "other"
  | "other insurance"
  | "other annuity"
  | "prif"
  | "rdsp"
  | "resp"
  | "rlif"
  | "rrif"
  | "pension"
  | "profit sharing plan"
  | "retirement"
  | "roth"
  | "roth 401k"
  | "rrsp"
  | "sep ira"
  | "simple ira"
  | "sipp"
  | "stock plan"
  | "thrift savings plan"
  | "tfsa"
  | "trust"
  | "ugma"
  | "utma"
  | "variable annuity"
  | "credit card"
  | "paypal"
  | "cd"
  | "checking"
  | "savings"
  | "money market"
  | "prepaid"
  | "auto"
  | "business"
  | "commercial"
  | "construction"
  | "consumer"
  | "home equity"
  | "loan"
  | "mortgage"
  | "overdraft"
  | "line of credit"
  | "student"
  | "cash management"
  | "keogh"
  | "mutual fund"
  | "recurring"
  | "rewards"
  | "safe deposit"
  | "sarsep"
  | "payroll"
  | "null";
