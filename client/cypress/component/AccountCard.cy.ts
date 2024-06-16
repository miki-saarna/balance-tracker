import AccountCard from "../../src/components/AccountCard.vue";

const exampleAccount = {
  account_id: "EXawBaoqEdS54z7e9X11uX8ywVJ3JwCr3D1Kq",
  balances: { available: 100 },
  name: "Plaid Checking",
  subtype: "checking",
};

describe("AccountCard component", () => {
  before(() => {
    cy.window().then((win) => {
      const tooltipContainer = win.document.createElement("div");
      tooltipContainer.id = "tooltip-container";
      win.document.body.appendChild(tooltipContainer);
    });
  });

  beforeEach(() => {
    cy.mount(AccountCard, {
      props: {
        name: exampleAccount.name,
        type: exampleAccount.subtype,
        balance: exampleAccount.balances.available,
      },
    });
  });

  it("properly renders account details", () => {
    cy.get("div").should("contains.text", exampleAccount.name);
    cy.get("div").should("contains.text", exampleAccount.subtype);
    cy.get("div").should(
      "contains.text",
      `$${exampleAccount.balances.available}`
    );
  });

  it("opens tooltip when toggled", () => {
    cy.get("#tooltip-container").children().should("have.length", 0);
    const tooltipToggler = cy.get("button");
    tooltipToggler.click();
    cy.get("#tooltip-container").children().should("have.length", 1);
    tooltipToggler.click();
    cy.get("#tooltip-container").children().should("have.length", 0);
  });
});
