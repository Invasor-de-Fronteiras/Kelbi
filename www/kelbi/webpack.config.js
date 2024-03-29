const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ES3Plugin = require('webpack-es3-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

const devMode = process.env.NODE_ENV !== 'production';

const productionPlugins = [new MiniCssExtractPlugin(), new ES3Plugin()];

const outputPath = path.resolve(__dirname, devMode ? './dev-build' : './build');

module.exports = {
  entry: ['@babel/polyfill', './src/index.tsx'],
  output: { path: outputPath, filename: '[name].js', chunkFilename: '[id].[chunkhash].js' },
  // mode: process.env.NODE_ENV || "development",
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
  },
  devServer: { static: outputPath },
  optimization: {
    runtimeChunk: 'single',
  },
  module: {
    rules: [
      // {
      //   test: /\.(js|jsx)$/,
      //   exclude: /node_modules/,
      //   use: ['babel-loader'],
      // },
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: ['babel-loader'],
      },
      {
        test: /\.(css)$/,
        use: [
          devMode ? 'style-loader' : MiniCssExtractPlugin.loader,
          'css-loader',
          'css-modules-typescript-loader',
          'postcss-loader',
        ],
      },
      {
        test: /\.(jpg|jpeg|png|gif|mp3|svg)$/,
        use: ['file-loader'],
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.join(__dirname, 'public', 'index.html'),
    }),
    ...(devMode ? [] : productionPlugins),
  ],
};
