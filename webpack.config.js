'use strict';

const process = require('process');
const webpack = require('webpack');
const path = require('path');
const buildPath = path.resolve(__dirname, 'static');
const nodeModulesPath = path.resolve(__dirname, 'node_modules');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const TransferWebpackPlugin = require('transfer-webpack-plugin'); // dev-server only

const version_hash = process.env.VERSION_HASH ? process.env.VERSION_HASH : '_version_';

var devtool = 'source-map'; // Render source-map file for final build
var plugins = [
  new CopyWebpackPlugin([
    { context: 'frontend', from: '**/*.css' },
    { context: 'frontend', from: '**/*.html' },
    { context: 'frontend', from: '**/*.png' },
    { context: 'frontend', from: '**/*.jpg' },
    { context: 'frontend', from: '**/*.eot' },
    { context: 'frontend', from: '**/*.svg' },
    { context: 'frontend', from: '**/*.ttf' },
    { context: 'frontend', from: '**/*.woff' },
    { context: 'frontend', from: '**/*.woff2' },
    { context: nodeModulesPath, from: 'bootstrap/dist/**', to: 'lib' },
    { context: nodeModulesPath, from: 'bootswatch/dist/yeti/**', to: 'lib' },
    { context: nodeModulesPath, from: 'jquery/dist/**', to: 'lib' },
    { context: nodeModulesPath, from: 'popper.js/dist/popper.*', to: 'lib'},
    { context: nodeModulesPath, from: 'react/umd/react.*', to: 'lib' },
    { context: nodeModulesPath, from: 'react-dom/umd/react-dom.*', to: 'lib' },
    { context: nodeModulesPath, from: '@fortawesome/fontawesome-free', to: 'lib/fontawesome-free' },
  ])
];

if (process.env.WEBPACK_DEV_SERVER === 'true') {
  console.log('WEBPACK_DEV_SERVER is true');
  devtool = 'eval';
  plugins.push(
    new TransferWebpackPlugin([
      {from: 'css'},
      {from: 'html'},
      {from: 'image'},
      {from: nodeModulesPath + '/bootstrap', to: 'lib'},
      {from: nodeModulesPath + '/bootswatch', to: 'lib'},
      {from: nodeModulesPath + '/jquery', to: 'lib'},
      {from: nodeModulesPath + '/popper.js', to: 'lib'},
      {from: nodeModulesPath + '/react', to: 'lib'},
      {from: nodeModulesPath + '/react-dom', to: 'lib'},
      {from: nodeModulesPath + '/@fortawesome/fontawesome-free', to: 'lib/fontawesome-free' },
    ], path.resolve(__dirname, "frontend"))
  );
}

const config = {
  mode: process.env.MINIFY === 'true' ? 'production' : 'development',
  entry: {
    main: './frontend/js/main.js',
    // setting: './frontend/js/setting.js',
  },
  resolve: {
    //When require, do not have to add these extensions to file's name
    extensions: ['.js', '.jsx', '.json', '.css'],
    //node_modules: ["web_modules", "node_modules"]  (Default Settings)
  },
  output: {
    path: path.join(buildPath, version_hash),
    publicPath: path.join('/static', version_hash),
    filename: "js/[name].bundle.js",
    chunkFilename: "js/[name].chunk.js"
  },
  externals: {
    'bootstrap': 'bootstrap',
    'bootswatch': 'bootswatch',
    'jquery': 'jQuery',
    'popper.js': 'popper.js',
    'react': 'React',
    'react-dom': 'ReactDOM',
  },
  optimization: {
    runtimeChunk: {
      name: 'vendor'
    },
    splitChunks: {
      name: 'vendor',
      chunks: 'initial',
    }
  },
  plugins: plugins,
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        //include: paths.appSrc,
        loader: require.resolve('babel-loader'),
        exclude: /node_modules/,
        options: {
          // This is a feature of `babel-loader` for Webpack (not Babel itself).
          // It enables caching results in ./node_modules/.cache/babel-loader/
          // directory for faster rebuilds.
          cacheDirectory: true,
          //plugins: ['react-hot-loader/babel'],
          presets: [
            ['@babel/react'],
            ['@babel/env', {
                "targets": {
                  "browsers": ["last 2 versions", "safari >= 7"]
                }
            }]
          ]
        }
      },
      {
        test: /\.css$/,
        loader: "style-loader!css-loader"
      },
      {
        test: /\.png$/,
        loader: "url-loader?limit=100000"
      },
      {
        test: /\.jpg$/,
        loader: "file-loader"
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url?limit=10000&mimetype=application/font-woff'
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url?limit=10000&mimetype=application/octet-stream'
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file'
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url?limit=10000&mimetype=image/svg+xml'
      }
    ],
  },
  // Dev server Configuration options
  devServer: {
    contentBase: 'frontend',  // Relative directory for base of server
    hot: true,        // Live-reload
    inline: true,
    port: 5000,
    host: '0.0.0.0',  // Change to '0.0.0.0' for external facing server
    proxy: {
      '*': {
        target: 'http://localhost:5001',
        secure: false,
        // bypass: function (req, res, proxyOptions) {
        //     const accept = req.headers.accept
        //     console.log(accept);
        //     //if (accept.indexOf('html') !== -1 || accept.indexOf('js') !== -1 || accept.indexOf('css') !== -1) {
        //         console.log('Skipping proxy for browser request.');
        //         return false;
        //     //}
        // }
      }
    },
  },
};

module.exports = config;
