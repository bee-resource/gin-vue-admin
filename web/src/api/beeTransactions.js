import service from "@/utils/request";

export const getBeeTransactionList = (params) => {
  return service({
    url: "/beeTransactions/getBeeTransactionList",
    method: "get",
    params,
  });
};
