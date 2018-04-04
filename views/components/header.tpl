{{/* Header of every web page, needs root data (pass "." as data) */}}
<div class="header">
  <div class="container">
    <div class="title-bar">
      <div>
          {{if .IsLogin}}
          <div class="button-group">
              <a href="{{urlfor "SubmitController.Submit"}}{{if .Topic}}?topic={{.Topic.Name}}{{end}}">
                <button class="text-button">Submit {{if .Topic}}to {{.Topic.Name}}{{end}}</button>
              </a>
              <a href="{{urlfor "SubmitController.CreateTopic"}}">
                <button class="text-button">Create Topic</button>
              </a>
            </div>
          {{end}}
      </div>
      <div>
        <a href="/" class="no-style">
          <h1 class="title">social-loot</h1>
        </a>
      </div>
      <div>
        <div class="button-group">
          <a href="/search">
            <button class="icon-button ion-ios-search"></button>
          </a>
          {{if .IsLogin}}
          <button class="icon-button large ion-ios-person-outline"
                  onclick="toggleUserPopup(this)">
          </button>
          {{else}}
          <a href="{{urlfor "LoginController.Login"}}">
            <button class="text-button">Login</button>
          </a>
          <a href="{{urlfor "LoginController.Signup"}}">
            <button class="text-button">Register</button>
          </a>
          {{end}}
        </div>
      </div>
    </div>
    {{if .Topics}} {{template "components/topic_list.tpl" .}} {{end}}
  </div>
  {{if .IsLogin}}
  <div class="popup-list" id="user-popup">
    <ul class="inner">
      <li>
        <a href="{{URL .User}}">
          Profile
        </a>
      </li>
      <li>
        <form action="{{urlfor "LoginController.Logout"}}" method="get" id="logout" style="display:inline">
          <input type="hidden" name="dest" value="/" />
          <a href="javascript:void(0)" onclick="document.getElementById('logout').submit()">Logout</a>
        </form>
      </li>
    </ul>
  </div>
  {{end}}
</div>
<!--
<div>
  <a href="/">Home</a>
  {{if .IsLogin}}
    <a itemprop="url" href='{{urlfor "SubmitController.Submit"}}{{if .Topic}}?topic={{.Topic.Name}}{{end}}'>Submit {{if .Topic}}to {{.Topic.Name}}{{end}}</a>
    <a itemprop="url" href='{{urlfor "SubmitController.CreateTopic"}}'>Create Topic</a>
  {{end}}
  <div class="right">
  {{if .IsLogin}}
    <span>User: {{.User.Name}}</span>
    <form action="{{urlfor "LoginController.Logout"}}" method="get" id="logout" style="display:inline">
      <input type="hidden" name="dest" value="/" />
      <a href="javascript:void(0)" onclick="document.getElementById('logout').submit()">(logout)</a>
    </form>
  {{else}}
    <a itemprop="url" href='{{urlfor "LoginController.Login"}}?dest={{.URL}}'>Login</a>
    <a itemprop="url" href='{{urlfor "LoginController.Signup"}}'>Signup</a>
  {{end}}
  </div>
</div> 
-->