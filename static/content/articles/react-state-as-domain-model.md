---
title: ドメインのstateを状態遷移として表す
image: images/Microsoft-Fluentui-Emoji-3d-Counterclockwise-Arrows-Button-3d.1024.png
publishedAt: 2026-06-28
updatedAt: 2026-06-28
---
React公式ドキュメントの「[stateの管理](https://ja.react.dev/learn/managing-state)」を読んだことがある方なら、UIを宣言的に設計するために `status` という状態（state）を導入する手法に見覚えがあるでしょう。

あるWebサイトでは、フォームの挙動を制御するために以下のようなステートマシンのアプローチが紹介されています。
```typescript
type FormStatus = 'typing' | 'submitting' | 'success' | 'error';
```
これはUIの表示制御において非常に強力なアプローチですが、一歩進めて「ビジネスルール（ドメインモデル）」の観点からこの`status`を捉え直すと、フロントエンドの堅牢性はさらに向上します。

本記事では、この`status`が必要なドメインモデルを題材に、開発者が書きがちなコードと、状態遷移の思想を取り入れた適切なコードを対比しながら解説します。

## テーマ: オンライン試験の解答提出
今回は、公式ドキュメントのフォームの例を発展させ、以下のビジネスルールを持つ「オンライン試験の解答提出」というドメインを考えます。

### ビジネスルール（不変条件）

1. 解答中 (`typing`): ユーザーは自由に解答を編集できる。
2. 提出中 (`submitting`): サーバーへの送信中。二重送信を防ぐため、解答の編集や再提出ボタンの連打は不可。
    
3. 提出成功 (`success`): 無事に受理された状態。一度成功したら、二度と解答の編集も再提出もできない。
    
4. エラー (`error`): 送信失敗。エラーメッセージを表示し、再試行できる。
## 1. 開発者が書きがちなコード
まずは、多くの現場で見かける、UIの都合だけで状態をバラバラに管理してしまっているコードです。
```tsx
import React, { useState } from "react";

export const BadExamSubmission: React.FC = () => {
  const [answer, setAnswer] = useState("");
  const [status, setStatus] = useState<"typing" | "submitting" | "success" | "error">("typing");
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleSubmit = async () => {
    // 連打防止のためのチェックがUI層に漏れ出ている
    if (status === "submitting") return;
    
    // 「一度成功したら提出できない」というルールが画面側の防御に依存している
    if (status === "success") {
      alert("既に提出済みです");
      return;
    }

    if (answer.trim() === "") {
      alert("解答が空です");
      return;
    }

    setStatus("submitting");
    try {
      // 擬似APIコール
      await saveAnswer(answer);
      setStatus("success");
    } catch (e) {
      setStatus("error");
      setErrorMessage("送信に失敗しました。内容を確認して再試行してください。");
    }
  };

  return (
    <div>
      {/* statusがsuccessの時でも、コードの書き方次第でtextareaが活性化してしまうリスクがある */}
      <textarea
        value={answer}
        onChange={(e) => setAnswer(e.target.value)}
        disabled={status === "submitting" || status === "success"}
      />
      
      <button onClick={handleSubmit} disabled={status === "submitting"}>
        {status === "submitting" ? "提出中..." : "試験を提出する"}
      </button>

      {status === "error" && <p style={{ color: "red" }}>{errorMessage}</p>}
      {status === "success" && <p style={{ color: "green" }}>提出が完了しました。</p>}
    </div>
  );
};

// 擬似的なAPI関数
const saveAnswer = (ans: string) => new Promise((res) => setTimeout(res, 1000));
```
### このコードの何が悪いのか？
一見正しく動くように見えますが、以下の致命的な課題を抱えています。
- `setStatus("success")`や`setStatus("typing")`をどこからでも呼び出せるため、プログラムのバグによって「提出成功(`success`)から解答中(`typing`)に逆戻りする」といった不正な状態遷移を簡単に引き起こせます。
- 「一度提出したら変更できない」「空の解答は出せない」といったドメインのルールが、コンポーネントのイベントハンドラーの中に`if`文として埋もれています。
## 2. 適切なコード
次に、「stateはドメインモデルの振る舞い（状態遷移）そのものである」という思想に基づいた設計です。

TypeScriptの「判別可能なユニオン型（Discriminated Unions）」を使い、その状態の時にしか存在し得ないデータを型で縛り、状態遷移のルールを純粋関数（ドメイン層）に隔離します。
### 1. ドメイン層
Reactのことを一切知らない、純粋なビジネスルールのみを持つファイルです。
```typescript
// 状態ごとのデータ構造を厳格に定義
export type TypingState = { status: "typing"; answer: string };
export type SubmittingState = { status: "submitting"; answer: string };
export type SuccessState = { status: "success"; finalizedAnswer: string; submittedAt: Date };
export type ErrorState = { status: "error"; answer: string; errorReason: string };

// stateの定義
export type ExamState = TypingState | SubmittingState | SuccessState | ErrorState;

// 起きうるイベントの定義
export type ExamEvent =
  | { type: "EDIT"; payload: { text: string } }
  | { type: "START_SUBMIT" }
  | { type: "SUBMIT_SUCCESS"; payload: { submittedAt: Date } }
  | { type: "SUBMIT_FAILURE"; payload: { reason: string } };

// 状態遷移のルール（ドメインモデルの振る舞い）
export const ExamDomain = {
  // 解答の編集はtypingやerrorの時だけ受け付ける
  edit(state: TypingState | ErrorState, text: string): TypingState {
    return { status: "typing", answer: text };
  },

  // 提出開始（typingやerrorからのみ遷移可能）
  startSubmit(state: TypingState | ErrorState): SubmittingState | ErrorState {
    if (state.answer.trim() === "") {
      return {
        status: "error",
        answer: state.answer,
        errorReason: "解答が空の状態で提出することはできません"
      };
    }
    return { status: "submitting", answer: state.answer };
  },

  // 提出成功（submittingからのみ遷移可能。成功データに変換され、これ以降編集不可に）
  confirmSuccess(state: SubmittingState, submittedAt: Date): SuccessState {
    return {
      status: "success",
      finalizedAnswer: state.answer,
      submittedAt
    };
  },

  // 提出失敗
  handleFailure(state: SubmittingState, reason: string): ErrorState {
    return {
      status: "error",
      answer: state.answer,
      errorReason: reason
    };
  }
};
```
### 2. プレゼンテーション層
UI側では、`useReducer`を使ってこのドメインロジックを呼び出します。Reducerは自ら判断せず、ドメインモデルに「今の状態を次の状態に変えてくれ」と委ねるだけの存在になります。
```tsx
import React, { useReducer, useEffect } from "react";
import { ExamState, ExamEvent, ExamDomain } from "../domain/ExamDomain";

const examReducer = (state: ExamState, event: ExamEvent): ExamState => {
  // 現在の状態（status）に応じて、ドメインモデルの対応する振る舞いのみを呼び出す
  switch (state.status) {
    case "typing":
    case "error":
      if (event.type === "EDIT") return ExamDomain.edit(state, event.payload.text);
      if (event.type === "START_SUBMIT") return ExamDomain.startSubmit(state);
      return state;

    case "submitting":
      if (event.type === "SUBMIT_SUCCESS") return ExamDomain.confirmSuccess(state, event.payload.submittedAt);
      if (event.type === "SUBMIT_FAILURE") return ExamDomain.handleFailure(state, event.payload.reason);
      return state;

    case "success":
      // success状態はどのようなイベントが来ても状態遷移を拒否
      return state;

    default:
      return state;
  }
};

const initialState: ExamState = { status: "typing", answer: "" };

export const GoodExamSubmission: React.FC = () => {
  const [state, dispatch] = useReducer(examReducer, initialState);

  useEffect(() => {
    if (state.status !== "submitting") return;

    const executeSubmission = async () => {
      try {
        // 実際の送信処理（state.answerは型ガードにより安全に取得可能）
        await saveAnswer(state.answer); 
        
        dispatch({ type: "SUBMIT_SUCCESS", payload: { submittedAt: new Date() } });
      } catch (error) {
        dispatch({ type: "SUBMIT_FAILURE", payload: { reason: "通信エラーが発生しました" } });
      }
    };

    executeSubmission();
  }, [state.status, state.status === "submitting" ? state.answer : undefined]);

  const handleSubmit = () => {
    dispatch({ type: "START_SUBMIT" });
  };

  return (
    <div>
      {/* 型に基づいた安全なUIの出し分け */}
      {state.status !== "success" ? (
        <textarea
          value={state.answer}
          onChange={(e) => dispatch({ type: "EDIT", payload: { text: e.target.value } })}
          disabled={state.status === "submitting"}
        />
      ) : (
        <p>解答: <strong>{state.finalizedAnswer}</strong></p>
      )}

      {state.status !== "success" && (
        <button onClick={handleSubmit} disabled={state.status === "submitting"}>
          {state.status === "submitting" ? "提出中..." : "試験を提出する"}
        </button>
      )}

      {state.status === "error" && <p style={{ color: "red" }}>{state.errorReason}</p>}
      {state.status === "success" && (
        <p style={{ color: "green" }}>
          {state.submittedAt.toLocaleTimeString()} に提出が完了しました。
        </p>
      )}
    </div>
  );
};

const saveAnswer = (ans: string) => new Promise((res) => setTimeout(res, 1000));
```
## 何が変わったのか？
適切なコードへとリファクタリングしたことで、システムは以下のように改善されました。
### 不正な状態を表現できなくなった
`SuccessState` 型の定義を見ると、そこには`finalizedAnswer`（確定された文字）しか存在しません。また、`success`になった後はreducerがあらゆる遷移をブロックします。これにより、「提出後に勝手に値が書き換わるバグ」をコードの構造レベルで防げるようになりました。
### UI制御とビジネスルールが分離された
コンポーネント（UI層）の役割は「ボタンが押されたらイベントを通知する（`dispatch`）」ことだけです。「本当にそのタイミングで提出してよいか」「解答が空ではないか」というドメインのバリデーションはすべて`ExamDomain.ts`に集約されています。
### テスト可能性が向上した
「解答が空なら提出できない」「提出成功後は状態が変わらない」という仕様のテストを書く際、Reactのコンポーネントをレンダリングする必要はありません。ピュアなTypeScript関数である`ExamDomain.startSubmit`を呼び出すだけの、シンプルな単体テストを書くことで動作確認ができます。
## まとめ
React公式ドキュメントにかかれている`status`によるUI管理は良いアプローチです。しかし、それをコンポーネントの中だけで完結させるのではなく、「ドメインモデルが持つべき状態遷移のルール」として外側に切り出すことで、フロントエンドの複雑性が下がります。

stateをただの「データの器」にせず、「ビジネスルールを表現するステートマシン」としてモデリングすることで、「UIの都合」と「ビジネスの都合」がきれいに分離され、開発者の認知負荷も下がるという利点もあります。

これからのフロントエンドには「複雑化するWebアプリのルールをどう制御するか」という、アーキテクチャの視点が不可欠です。