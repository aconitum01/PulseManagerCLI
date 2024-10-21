# PulseManagerCLI

**PulseManagerCLI** は、コマンドラインでオーディオ出力デバイスを簡単に選択・管理するためのツールです。このツールは、複数のオーディオ出力デバイスを切り替えたり、ボリュームを調整するためのインタラクティブなインターフェースを提供します。

## 特徴
- 複数のオーディオ出力デバイスのリストを表示し、選択を可能にします。
- 簡単にデフォルトのオーディオ出力デバイスを設定できます。
- ボリュームをリアルタイムで調整できます。
- キーボード操作により、使いやすく設計されています。

## システム要件
- **PulseAudio** がインストールされ、オーディオ管理に使用されている Linux 環境

## インストール
以下の手順で **PulseManagerCLI** をインストールできます。

1. このリポジトリをクローンします。
   ```sh
   git clone https://github.com/yourusername/PulseManagerCLI.git
   cd PulseManagerCLI
   ```

2. バイナリファイルをビルドします。
   ```sh
   go build -o pulsemanager main.go
   ```

   - ビルド後、`pulsemanager` という名前のバイナリが生成されます。

## 使用方法
`pulsemanager` コマンドを実行することで、オーディオ出力デバイスの管理ができます。

```sh
./pulsemanager
```

以下の操作が可能です:
- **上下の矢印キー** または **`j` / `k` キー**: オーディオデバイスの選択を移動します。
- **左右の矢印キー** または **`h` / `l` キー**: ボリュームの調整を行います。
- **数字キー**: 該当番号のオーディオデバイスを直接選択します。
- **Enterキー**: 選択したオーディオデバイスをデフォルトに設定します。
- **`q` キー** または **Escキー**: プログラムを終了します。

## 例
```sh
./pulsemanager
```
このコマンドを実行すると、現在接続されているオーディオ出力デバイスのリストが表示されます。矢印キーを使って選択したり、Enterキーでデフォルトデバイスを設定することができます。

## 注意事項
- **PulseManagerCLI** は PulseAudio を利用しているため、PulseAudio がインストールされていない環境では動作しません。

