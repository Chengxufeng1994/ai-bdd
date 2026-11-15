# 驗收旅程範例：不執行，沒有 step definition。
# 規則與例子來自 ../../spec/examples/minimal-spec.md 的 US-1（FR-1、AC-1.1–1.3）。

Feature: 訪客出場結帳
  訪客進場抽紙票，出場前在繳費機結帳，付完就能開走。
  前 30 分鐘免費，超過才收費。

  Rule: 前 30 分鐘免費，超過才收費

    @covers-1.1 @covers-1.2
    Example: 訪客進來辦點事，半小時內就走
      Given 一位訪客開車進場並抽了紙票
      When 訪客停了 29 分鐘後到繳費機結帳
      Then 繳費機顯示不用付錢
      When 另一位訪客停滿整整 30 分鐘才來結帳
      Then 繳費機一樣顯示不用付錢

    @covers-1.3
    Example: 訪客多停了一分鐘
      Given 一位訪客開車進場並抽了紙票
      When 訪客停了 31 分鐘後到繳費機結帳
      Then 繳費機顯示要付 30 元
