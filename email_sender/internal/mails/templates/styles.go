package templates

var StartTempStyle = `<html>
<head>
   <meta charset="UTF-8">
   <meta name="x-apple-disable-message-reformatting">
   <meta name="viewport" content="width=device-width,initial-scale=1">
   <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
   <title>Exchange Rates</title>
   <style>
   img {
	border: none;
	-ms-interpolation-mode: bicubic;
	max-width: 100%;
  }

  table {
	border-collapse: collapse;
	min-width: 100%;
	width: 100%;
  }

  table td {
	font-family: sans-serif;
	font-size: 14px;
	vertical-align: top;
  }

  .body {
	background-color: #ebf6ff;
	width: 100%;
  }

  .container {
	display: block;
	margin: 0 auto !important;
	max-width: 800px;
	padding: 10px;
	width: 800px;
  }

  .content {
	box-sizing: border-box;
	display: block;
	margin: 0 auto;
	max-width: 800px;
	padding: 10px;
  }

  .main {
	background: #ffffff;
	border-radius: 18px;
	width: 100%;
  }

  .header {
	padding: 20px 0;
  }

  .wrapper {
	box-sizing: border-box;
	padding: 20px;
  }

  .content-block {
	padding-bottom: 10px;
	padding-top: 10px;
  }

  .footer {
	clear: both;
	margin-top: 10px;
	text-align: center;
	width: 100%;
  }

  .footer td,
  .footer p,
  .footer span,
  .footer a {
	color: #9a9ea6;
	font-size: 12px;
	text-align: center;
  }

  h1,
  h2,
  h3,
  h4 {
	color: #06090f;
	font-family: sans-serif;
	font-weight: 400;
	line-height: 1.4;
	margin: 0;
	margin-bottom: 30px;
  }

  h1 {
	font-size: 35px;
	font-weight: 300;
	text-align: center;
	text-transform: capitalize;
  }

  p,
  ul,
  ol {
	font-family: sans-serif;
	font-size: 14px;
	font-weight: normal;
	margin: 0;
	margin-bottom: 15px;
  }

  p li,
  ul li,
  ol li {
	list-style-position: inside;
	margin-left: 5px;
  }

  a {
	color: #00457c;
	text-decoration: underline;
  }

  .btn {
	box-sizing: border-box;
	min-width: 100%;
	width: 100%;
  }

  .btn>tbody>tr>td {
	padding-bottom: 15px;
  }

  .btn table {
	min-width: auto;
	width: auto;
  }

  .btn table td {
	background-color: #ffffff;
	border-radius: 5px;
	text-align: center;
  }

  .btn a {
	background-color: #ffffff;
	border: solid 1px #00457c;
	border-radius: 5px;
	box-sizing: border-box;
	color: #00457c;
	cursor: pointer;
	display: inline-block;
	font-size: 14px;
	font-weight: bold;
	margin: 0;
	padding: 12px 25px;
	text-decoration: none;
	text-transform: capitalize;
  }

  .btn-primary table td {
	background-color: #00457c;
  }

  .btn-primary a {
	background-color: #00457c;
	border-color: #00457c;
	color: #ffffff;
  }

  .last {
	margin-bottom: 0;
  }

  .first {
	margin-top: 0;
  }

  .align-center {
	text-align: center;
  }

  .align-right {
	text-align: right;
  }

  .align-left {
	text-align: left;
  }

  .clear {
	clear: both;
  }

  .mt0 {
	margin-top: 0;
  }

  .mb0 {
	margin-bottom: 0;
  }

  .preheader {
	color: transparent;
	display: none;
	height: 0;
	max-height: 0;
	max-width: 0;
	opacity: 0;
	overflow: hidden;
	visibility: hidden;
	width: 0;
  }

  .powered-by a {
	text-decoration: none;
  }

  hr {
	border: 0;
	border-bottom: 1px solid #f6f6f6;
	margin: 20px 0;
  }

  @media only screen and (max-width: 620px) {
	table[class="body"] h1 {
	  font-size: 28px !important;
	  margin-bottom: 10px !important;
	}

	table[class="body"] p,
	table[class="body"] ul,
	table[class="body"] ol,
	table[class="body"] td,
	table[class="body"] span,
	table[class="body"] a {
	  font-size: 16px !important;
	}

	table[class="body"] .wrapper,
	table[class="body"] .article {
	  padding: 10px !important;
	}

	table[class="body"] .content {
	  padding: 0 !important;
	}

	table[class="body"] .container {
	  padding: 0 !important;
	  width: 100% !important;
	}

	table[class="body"] .main {
	  border-left-width: 0 !important;
	  border-radius: 0 !important;
	  border-right-width: 0 !important;
	}

	table[class="body"] .btn table {
	  width: 100% !important;
	}

	table[class="body"] .btn a {
	  width: 100% !important;
	}

	table[class="body"] .img-responsive {
	  height: auto !important;
	  max-width: 100% !important;
	  width: auto !important;
	}
  }

  @media all {
	.ExternalClass {
	  width: 100%;
	}

	.ExternalClass,
	.ExternalClass p,
	.ExternalClass span,
	.ExternalClass font,
	.ExternalClass td,
	.ExternalClass div {
	  line-height: 100%;
	}

	.apple-link a {
	  color: inherit !important;
	  font-family: inherit !important;
	  font-size: inherit !important;
	  font-weight: inherit !important;
	  line-height: inherit !important;
	  text-decoration: none !important;
	}

	.btn-primary table td:hover {
	  background-color: #063b67 !important;
	}

	.btn-primary a:hover {
	  background-color: #063b67 !important;
	  border-color: #063b67 !important;
	}
  }

  .table-bordered{
	border-collapse: collapse;
  }

  .table-bordered table thead {
	background: #00457c;
	color: #fff;
	border: 1px solid #00457c;
  }

  .table-bordered table thead tr th{
	padding: 8px;
  }

 .table-bordered table td {
	border: 1px solid #9a9ea6;
	padding: 8px;
  }
</style>
</head>
<body
   style="background-color: #ebf6ff;font-family: sans-serif; -webkit-font-smoothing: antialiased; font-size: 14px; line-height: 1.4; margin: 0; padding: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;">
   <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body" style="font-family: sans-serif; font-size: 14px; vertical-align: top;">
	  <tr>
		 <td>&nbsp;</td>
		 <td class="container" style="display: block; margin: 0 auto !important; max-width: 580px; padding: 10px; width: 580px;">
			<div class="header" style="padding: 20px 0;">
			   <table role="presentation" border="0" cellpadding="0" cellspacing="0" style="border-collapse: separate; min-width: 100%; width: 100%;">
				  <tr>
					 <td class="align-center" style="text-align: center;">
						<a ><img src="https://extension.umd.edu/sites/extension.umd.edu/files/styles/optimized/public/2023-01/10-IRS-audit-triggers.jpg?itok=7jmgHjrs" width="200"
						   height="60" alt="ExchangeRates" style="border: none;  -ms-interpolation-mode: bicubic; max-width: 100%;" /></a>
					 </td>
				  </tr>
			   </table>
			</div>
			<div class="content">
			   <span class="preheader" style="color: transparent; display: none; height: 0; max-height: 0; max-width: 0; opacity: 0; overflow: hidden; visibility: hidden; width: 0;"></span>
			   <table role="presentation" class="main" style="border-collapse: separate; min-width: 100%; width: 100%; background: #ffffff; border-radius: 18px; width: 100%;">
				  <tr>
					 <td class="wrapper" style="box-sizing: border-box; padding: 20px;">
						<table role="presentation" border="0" cellpadding="0" cellspacing="0"  style="border-collapse: separate; min-width: 100%; width: 100%;">
						   <tr>
							  <td>`
var EndTempStyle = ` </td>
</tr>
</table>
</td>
</tr>
</table>
</td>
</tr>
</table>
<div class="footer" style="clear: both; margin-top: 10px; text-align: center; width: 100%;">
<table role="presentation" border="0" cellpadding="0" cellspacing="0"  style="border-collapse: separate; min-width: 100%; width: 100%;">
<tr>
<td class="content-block powered-by" style="padding-bottom: 10px; padding-top: 10px; color: #9a9ea6; font-size: 12px; text-align: center;">
Powered by <a style="color: #9a9ea6; font-size: 12px; text-align: center;" >Exchange Rates Team</a>.
</td>
</tr>
</table>
</div>
</div>
</td>
<td>&nbsp;</td>
</tr>
</table>
</body>
</html>`
