<template>
	<div>
		<button class="button is-small" :class="color" @click="show" @click.stop>{{message}}</button>
		<div class="hiddens" v-if="active">
			<div class="ant-popover ant-popover-placement-top" @click.stop :class="active ? '' : 'ant-popover-hidden'" :style="active ?`top: ${top}px; left: ${left}px; transform-origin: 50% 103px;`:''">
				<div class="ant-popover-content">
					<div class="ant-popover-arrow">
					</div>
					<div class="ant-popover-inner" role="tooltip">
						<div>
							<div class="ant-popover-inner-content">
								<div class="ant-popover-message">
									<i aria-label="icon: exclamation-circle" class="anticon anticon-exclamation-circle">
										<svg viewBox="64 64 896 896" focusable="false" class="" data-icon="exclamation-circle"
										width="1em" height="1em" fill="currentColor" aria-hidden="true">
											<path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm-32 232c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V296zm32 440a48.01 48.01 0 0 1 0-96 48.01 48.01 0 0 1 0 96z">
											</path>
										</svg>
									</i>
									<div class="ant-popover-message-title">
										确定要{{message}}吗？
									</div>
								</div>
								<div class="ant-popover-buttons">
									<button type="button" class="ant-btn ant-btn-sm" @click="hidden">
										<span>
											否
										</span>
									</button>
									<button type="button"  @click="submit" class="ant-btn ant-btn-primary ant-btn-sm">
										<span>
											是
										</span>
									</button>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
import { reactive, toRefs, onMounted, onUnmounted, defineComponent } from 'vue'

export default defineComponent({
  name: "PopoButton",
	props: {
		message: String,
		color: String,
		callBack: Function
	},
  setup(props) {
		let states = reactive({
      active: false,
			left: 0,
			top: 0,
			hidden: Function
    })
		onMounted(() => {
      window.addEventListener('click', states.hidden)
    })
 
    onUnmounted(() => {
      window.removeEventListener('click', states.hidden)
    })
		const show = (e) => {
			const t = e.target
			states.left = t.offsetLeft - t.clientWidth - 6
			states.top = t.offsetTop - (t.clientHeight * 3) - 6
			states.active = !states.active
		}
		const hidden =() =>{
			states.active = false
		}
		states.hidden = hidden
		const submit =() => {
			states.active = false
			props.callBack()
		}

		return {
			...toRefs(states),
			show,
			hidden,
			submit
		}
  },
})
</script>


<style scoped>
.ant-popover {
	-webkit-box-sizing:border-box;
	box-sizing:border-box;
	margin:0;
	padding:0;
	color:rgba(0,0,0,.65);
	font-size:14px;
	font-variant:tabular-nums;
	line-height:1.5;
	list-style:none;
	-webkit-font-feature-settings:"tnum","tnum";
	font-feature-settings:"tnum","tnum";
	position:absolute;
	top:0;
	left:0;
	z-index:1030;
	font-weight:400;
	white-space:normal;
	text-align:left;
	cursor:auto;
	-webkit-user-select:text;
	-moz-user-select:text;
	-ms-user-select:text;
	user-select:text;
  width:160px
}
.ant-popover:after {
	position:absolute;
	background:hsla(0,0%,100%,.01);
	content:""
}
.ant-popover-hidden {
	display:none
}
.ant-popover-placement-top,.ant-popover-placement-topLeft,.ant-popover-placement-topRight {
	padding-bottom:10px
}
.ant-popover-placement-right,.ant-popover-placement-rightBottom,.ant-popover-placement-rightTop {
	padding-left:10px
}
.ant-popover-placement-bottom,.ant-popover-placement-bottomLeft,.ant-popover-placement-bottomRight {
	padding-top:10px
}
.ant-popover-placement-left,.ant-popover-placement-leftBottom,.ant-popover-placement-leftTop {
	padding-right:10px
}
.ant-popover-inner {
	background-color:#fff;
	background-clip:padding-box;
	border-radius:4px;
	-webkit-box-shadow:0 2px 8px rgba(0,0,0,.15);
	box-shadow:0 2px 8px rgba(0,0,0,.15);
	-webkit-box-shadow:0 0 8px rgba(0,0,0,.15)\9;
	box-shadow:0 0 8px rgba(0,0,0,.15)\9
}
@media (-ms-high-contrast:none),screen and (-ms-high-contrast:active) {
	.ant-popover-inner {
	-webkit-box-shadow:0 2px 8px rgba(0,0,0,.15);
	box-shadow:0 2px 8px rgba(0,0,0,.15)
}
}.ant-popover-title {
	min-width:177px;
	min-height:32px;
	margin:0;
	padding:5px 16px 4px;
	color:rgba(0,0,0,.85);
	font-weight:500;
	border-bottom:1px solid #e8e8e8
}
.ant-popover-inner-content {
	padding:12px 16px;
	color:rgba(0,0,0,.65)
}
.ant-popover-message {
	position:relative;
	padding:4px 0 12px;
	color:rgba(0,0,0,.65);
	font-size:14px
}
.ant-popover-message>.anticon {
	position:absolute;
	top:8px;
	color:#faad14;
	font-size:14px
}
.ant-popover-message-title {
	padding-left:22px
}
.ant-popover-buttons {
	margin-bottom:4px;
	text-align:right
}
.ant-popover-buttons button {
	margin-left:8px
}
.ant-popover-arrow {
	position:absolute;
	display:block;
	width:8.48528137px;
	height:8.48528137px;
	background:transparent;
	border-style:solid;
	border-width:4.24264069px;
	-webkit-transform:rotate(45deg);
	-ms-transform:rotate(45deg);
	transform:rotate(45deg)
}
.ant-popover-placement-top>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-topLeft>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-topRight>.ant-popover-content>.ant-popover-arrow {
	bottom:6.2px;
	border-color:transparent #fff #fff transparent;
	-webkit-box-shadow:3px 3px 7px rgba(0,0,0,.07);
	box-shadow:3px 3px 7px rgba(0,0,0,.07)
}
.ant-popover-placement-top>.ant-popover-content>.ant-popover-arrow {
	left:50%;
	-webkit-transform:translateX(-50%) rotate(45deg);
	-ms-transform:translateX(-50%) rotate(45deg);
	transform:translateX(-50%) rotate(45deg)
}
.ant-popover-placement-topLeft>.ant-popover-content>.ant-popover-arrow {
	left:16px
}
.ant-popover-placement-topRight>.ant-popover-content>.ant-popover-arrow {
	right:16px
}
.ant-popover-placement-right>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-rightBottom>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-rightTop>.ant-popover-content>.ant-popover-arrow {
	left:6px;
	border-color:transparent transparent #fff #fff;
	-webkit-box-shadow:-3px 3px 7px rgba(0,0,0,.07);
	box-shadow:-3px 3px 7px rgba(0,0,0,.07)
}
.ant-popover-placement-right>.ant-popover-content>.ant-popover-arrow {
	top:50%;
	-webkit-transform:translateY(-50%) rotate(45deg);
	-ms-transform:translateY(-50%) rotate(45deg);
	transform:translateY(-50%) rotate(45deg)
}
.ant-popover-placement-rightTop>.ant-popover-content>.ant-popover-arrow {
	top:12px
}
.ant-popover-placement-rightBottom>.ant-popover-content>.ant-popover-arrow {
	bottom:12px
}
.ant-popover-placement-bottom>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-bottomLeft>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-bottomRight>.ant-popover-content>.ant-popover-arrow {
	top:6px;
	border-color:#fff transparent transparent #fff;
	-webkit-box-shadow:-2px -2px 5px rgba(0,0,0,.06);
	box-shadow:-2px -2px 5px rgba(0,0,0,.06)
}
.ant-popover-placement-bottom>.ant-popover-content>.ant-popover-arrow {
	left:50%;
	-webkit-transform:translateX(-50%) rotate(45deg);
	-ms-transform:translateX(-50%) rotate(45deg);
	transform:translateX(-50%) rotate(45deg)
}
.ant-popover-placement-bottomLeft>.ant-popover-content>.ant-popover-arrow {
	left:16px
}
.ant-popover-placement-bottomRight>.ant-popover-content>.ant-popover-arrow {
	right:16px
}
.ant-popover-placement-left>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-leftBottom>.ant-popover-content>.ant-popover-arrow,.ant-popover-placement-leftTop>.ant-popover-content>.ant-popover-arrow {
	right:6px;
	border-color:#fff #fff transparent transparent;
	-webkit-box-shadow:3px -3px 7px rgba(0,0,0,.07);
	box-shadow:3px -3px 7px rgba(0,0,0,.07)
}
.ant-popover-placement-left>.ant-popover-content>.ant-popover-arrow {
	top:50%;
	-webkit-transform:translateY(-50%) rotate(45deg);
	-ms-transform:translateY(-50%) rotate(45deg);
	transform:translateY(-50%) rotate(45deg)
}
.ant-popover-placement-leftTop>.ant-popover-content>.ant-popover-arrow {
	top:12px
}
.ant-popover-placement-leftBottom>.ant-popover-content>.ant-popover-arrow {
	bottom:12px
}
.hiddens {
  position: absolute;
  width: 100%;
  height: 0px;
  left: 0px;
  top: 0px;
  right: 0px;
}
.ant-btn {
  line-height: 1.5;
  position: relative;
  display: inline-block;
  font-weight: 400;
  white-space: nowrap;
  text-align: center;
  background-image: none;
  -webkit-box-shadow: 0 2px 0 rgb(0 0 0 / 2%);
  box-shadow: 0 2px 0 rgb(0 0 0 / 2%);
  cursor: pointer;
  -webkit-transition: all .3s cubic-bezier(.645,.045,.355,1);
  transition: all .3s cubic-bezier(.645,.045,.355,1);
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  -ms-touch-action: manipulation;
  touch-action: manipulation;
  height: 32px;
  padding: 0 15px;
  font-size: 14px;
  border-radius: 4px;
  color: rgba(0,0,0,.65);
  background-color: #fff;
  border: 1px solid #d9d9d9;
}
.ant-btn-sm {
  height: 24px;
  padding: 0 7px;
  font-size: 14px;
  border-radius: 4px;
}
.ant-btn-primary {
  color: #fff;
  background-color: #1da57a;
  border-color: #1da57a;
  text-shadow: 0 -1px 0 rgb(0 0 0 / 12%);
  -webkit-box-shadow: 0 2px 0 rgb(0 0 0 / 5%);
  box-shadow: 0 2px 0 rgb(0 0 0 / 5%);
}
.button {
	margin-right: 0.5rem
}
</style>