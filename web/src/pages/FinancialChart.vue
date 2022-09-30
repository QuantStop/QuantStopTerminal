<template>
<div class="h-100" style="display: flex;">

    <!-- Chart -->
    <div class="chart-container">

      <!-- Top Toolbar -->
      <div class="chart-top-bar">

        <!-- Connection Light / Button -->
        <div class="connection-status">
          <StatusIndicator :value=dataConnected :height=15 :width=15></StatusIndicator>
        </div>

        <!-- Connection Status Message / Time -->
        <div class="icon-wrapper connection-time">
          <span v-if="!dataConnected" class="connection-message red">
            Disconnected <br>
            *unsupported
          </span>
          <span v-if="dataConnected" class="connection-message green">
            {{ this.date }} <br>
            {{ this.time }}
          </span>
        </div>

        <!-- Exchange Select -->
        <div class="icon-wrapper">
          <Dropdown
              :options="exchanges"
              v-on:selected="onExchangeSelected"
              v-on:filter="getExchangeValues"
              :disabled="false"
              placeholder="Please an exchange">
          </Dropdown>
        </div>

        <!-- Product Select -->
        <div class="icon-wrapper">
          <Dropdown
              :options="products"
              v-on:selected="onProductSelected"
              v-on:filter="getProductValues"
              :disabled="productInputDisabled"
              placeholder="Please a product">
          </Dropdown>
        </div>

        <!-- Periods -->
        <div class="icon-wrapper">
          <select id="chartPeriodSelect" class="form-control" v-model="category" @change="onPeriodChanged">
            <option v-for="template in templates"
                    :selected="template === category ? 'selected' : ''"
                    :value="template">
              {{ template }}
            </option>
          </select>
        </div>

        <!-- Indicators -->
        <div class="icon-wrapper dropdown" >
          <a type="button" id="dropdownMenuButton1" data-bs-toggle="dropdown" data-bs-auto-close="outside" aria-expanded="false" class="d-flex align-items-center">
            <svg viewBox="0 0 1024 1024">
              <path d="M768 921.6H275.2c-83.2 0-147.2-64-147.2-147.2V275.2C128 192 192 128 275.2 128H768c83.2 0 147.2 64 147.2 147.2 0 25.6-12.8 32-38.4 32s-38.4-12.8-38.4-32c0-38.4-32-70.4-70.4-70.4H275.2c-38.4 0-70.4 32-70.4 70.4v499.2c0 38.4 32 70.4 70.4 70.4H768c38.4 0 70.4-32 70.4-70.4 0-19.2 19.2-38.4 38.4-38.4s38.4 19.2 38.4 38.4c0 76.8-64 147.2-147.2 147.2z m-38.4-281.6c-6.4 0-19.2 0-25.6-6.4L582.4 499.2c-12.8-12.8-12.8-32 0-44.8 12.8-12.8 32-12.8 44.8 0l121.6 134.4c12.8 12.8 12.8 32 0 44.8 0 6.4-6.4 6.4-19.2 6.4z m-140.8 12.8c-6.4 0-12.8 0-19.2-6.4-12.8-12.8-12.8-25.6 0-38.4l166.4-172.8c12.8-12.8 25.6-12.8 38.4 0s12.8 25.6 0 38.4l-166.4 172.8c-6.4 6.4-12.8 6.4-19.2 6.4z m-44.8-377.6c12.8 0 44.8 6.4 44.8 38.4 0 25.6-38.4 19.2-51.2 19.2-57.6 0-76.8 32-83.2 57.6l-12.8 76.8h64c12.8 0 25.6 12.8 25.6 25.6s-12.8 25.6-25.6 25.6H435.2l-44.8 198.4s-6.4 57.6-44.8 57.6-19.2-57.6-19.2-57.6l44.8-192h-57.6c-12.8 0-25.6-12.8-25.6-25.6 0-25.6 12.8-32 25.6-32H384L403.2 384C416 275.2 512 268.8 544 275.2z m326.4 115.2c6.4-19.2 19.2-25.6 38.4-25.6 19.2 6.4 25.6 19.2 25.6 38.4l-64 256c-6.4 19.2-19.2 25.6-38.4 25.6-19.2-6.4-25.6-19.2-25.6-38.4l64-256z"></path>
            </svg>
            <span>Indicators</span>
          </a>
          <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1">
            <li><label class="list-group-item" v-for="type in mainTechnicalIndicatorTypes">
              <input
                  :key="type.name"
                  v-on:click="setCandleTechnicalIndicator(type)"
                  class="form-check-input me-1"
                  type="checkbox"
                  v-model="type.checked"
              >
              {{type.name}}
            </label></li>
            <li><label class="list-group-item">
              <input
                  v-on:click="setCandleTechnicalIndicator(customEmojiIndicatorTypes)"
                  class="form-check-input me-1"
                  type="checkbox"
                  v-model="customEmojiIndicatorTypes.checked"
              >
              Custom Indicator
            </label></li>

            <li><hr class="dropdown-divider"></li>

            <li><label class="list-group-item" v-for="type in subTechnicalIndicatorTypes">
              <input
                  :key="type.name"
                  v-on:click="setSubTechnicalIndicator(type)"
                  class="form-check-input me-1"
                  type="checkbox"
                  v-model="type.checked"
              >
              {{type.name}}
            </label></li>
          </ul>
        </div>

        <!-- Style -->
        <div class="icon-wrapper dropdown">
          <a type="button" id="dropdownMenuButton2" data-bs-toggle="dropdown" data-bs-auto-close="outside" aria-expanded="false" class="d-flex align-items-center">
            <svg viewBox="0 0 1024 1024" style="width: 22px; height: 22px;">
              <path d="M204.4 524.9c-14.5 1.5-26.2 13.2-27.7 27.7-2.1 19.9 14.6 36.7 34.6 34.6 14.5-1.5 26.2-13.2 27.8-27.8 2-19.9-14.8-36.6-34.7-34.5zM265.4 473.7c21.8-1.9 39.4-19.5 41.4-41.4 2.5-28.5-21.2-52.3-49.7-49.7-21.8 1.9-39.4 19.5-41.4 41.4-2.6 28.4 21.2 52.2 49.7 49.7zM415.8 266.9c-28.5 1.8-51.6 24.9-53.4 53.4-2.2 34.5 26.4 63.1 60.9 60.9 28.5-1.8 51.6-24.9 53.4-53.4 2.1-34.6-26.4-63.1-60.9-60.9zM621.9 253.8c-35.1 2.2-63.4 30.6-65.6 65.6-2.7 42.4 32.4 77.6 74.8 74.8 35.1-2.2 63.4-30.6 65.6-65.6 2.8-42.4-32.3-77.5-74.8-74.8zM966.5 276.4c-0.5-7.6-4-14.6-9.8-19.6l-0.7-0.6c-5.2-4.5-11.9-7-18.8-7-8.3 0-16.2 3.6-21.6 9.9L574 652.4l-43.5 85.5 1.1 0.9-4.9 11.3 11.1-5.9 1.5 1.3 78-54.3 342.3-394c5-5.8 7.4-13.2 6.9-20.8z"></path>
              <path d="M897.8 476.3c-13.8-1.4-26.7 7.4-30.4 20.7-6.9 24.6-19.3 64.5-35.1 97.8C809.5 643 767.4 710.1 696.7 756c-72.2 46.9-142.7 56.7-189.2 56.7-37 0-72.2-6.1-101.7-17.7-26.9-10.5-46.4-24.6-54.9-39.7-3.4-6.1-7.2-12.9-11.2-20.2-17.2-31.1-36.6-66.5-49.7-77.4-15.9-13.2-39.1-15-59.8-15-8.1 0-40.8 1.3-48.5 1.3-33.1 0-49.4-6.5-56.1-22.4-17.8-42.3-7.3-114.3 26.8-183.4C205.2 331.4 300 253.3 412.6 224c40-10.6 81.2-18.9 121.3-18.9 85.6 0 187.8 32.8 252.5 77.2 11.4 7.8 26.9 5.8 35.7-4.9 10.4-12.6 7.1-31.4-6.8-39.8-23.3-14-57.9-34-86.3-47.1-60.3-27.9-123.7-41.9-189.2-41.9-68.1 0-148.8 16.4-217.2 47.2-78.1 35-135.2 85-179.4 147.5-36.4 51.4-67.8 111.1-80.1 168.7-7.5 35.1-6.8 57.4-2.4 87.8 4.2 29.2 13.4 52.5 26.9 67.5 22.4 25.1 51.5 37.4 89 37.4 13.9 0 56.3-5 63.1-5 7.4 0 12.2 1.2 14.4 3.8 6.4 7.4 14.4 22.4 23.7 39.9 7.5 14.1 15.9 30.1 25.4 45.3 12.1 19.5 36.9 40.4 66.5 55.9 27 14.1 71.9 31 132.2 31 72 0 148.3-23.6 226.7-70.1 74.9-44.4 123-118.9 150.2-173.6 19-38.3 34.7-87.2 43.8-119.1 4.8-17.3-7-34.7-24.8-36.5z"></path>
            </svg>
            <span>Style</span>
          </a>
          <div class="dropdown-menu" aria-labelledby="dropdownMenuButton2">

            <div class="accordion accordion-flush" id="accordionFlushExample">

              <!-- Grid Styling -->
              <div class="accordion-item">
                <h2 class="accordion-header" id="flush-headingOne">
                  <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseOne" aria-expanded="false" aria-controls="flush-collapseOne">
                    Grid
                  </button>
                </h2>
                <div id="flush-collapseOne" class="accordion-collapse collapse" aria-labelledby="flush-headingOne" data-bs-parent="#accordionFlushExample">
                  <div class="accordion-body">
                    <div class="accordion-input d-flex justify-content-between">

                      <span>Show</span>
                      <ToggleSwitch
                          id="chartGridEnable"
                          :value="chartStyle.props.grid.show"
                          @input="setGridShow"
                          :width="50"
                          :height="25"
                      ></ToggleSwitch>

                    </div>
                    <div class="accordion-divider"></div>
                    <div class="accordion-item">
                      <h2 class="accordion-header" id="flush-headingOneOne">
                        <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseOneOne" aria-expanded="false" aria-controls="flush-collapseOneOne">
                          Horizontal Line
                        </button>
                      </h2>
                      <div id="flush-collapseOneOne" class="accordion-collapse collapse" aria-labelledby="flush-headingOneOne" data-bs-parent="#flush-collapseOne">
                        <div class="accordion-body">
                          <div class="accordion-input d-flex justify-content-between">
                            <span>Show</span>
                            <ToggleSwitch
                                id="chartGridEnable"
                                :value="chartStyle.props.grid.horizontal.show"
                                @input="setHorizontalGridShow"
                                :width="50"
                                :height="25"
                            ></ToggleSwitch>
                          </div>

                          <div class="accordion-input d-flex justify-content-between">
                            <span>Style</span>
                            <select
                                :value="chartStyle.props.grid.horizontal.style"
                                class="form-select form-select-sm"
                                aria-label=".form-select-sm"
                                @input="setGridHorizontalStyle"
                            >
                              <option value="dash">Dash</option>
                              <option value="solid">Solid</option>
                            </select>
                          </div>
                          <div class="accordion-input d-flex justify-content-between">
                            <span>Size</span>
                            <label for="hLineSize" class="form-label"></label>
                            <input type="range" class="form-range" min="1" max="10" step="1" id="hLineSize"
                                   @input="setGridHorizontalSize"
                                   :value="chartStyle.props.grid.horizontal.size"
                            >
                          </div>
                          <div class="accordion-input d-flex justify-content-between">
                            <span>Color</span>
                            <ColorPickerInput
                                class="w-50"
                                :color="chartStyle.props.grid.horizontal.color"
                                v-model="chartStyle.props.grid.horizontal.color"
                                @input="setGridHorizontalColor"
                            ></ColorPickerInput>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="accordion-item">
                      <h2 class="accordion-header" id="flush-headingOneTwo">
                        <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseOneTwo" aria-expanded="false" aria-controls="flush-collapseOneTwo">
                          Vertical
                        </button>
                      </h2>
                      <div id="flush-collapseOneTwo" class="accordion-collapse collapse" aria-labelledby="flush-headingOneTwo" data-bs-parent="#flush-collapseOne">
                        <div class="accordion-body">
                          <div class="accordion-input d-flex justify-content-between">
                            <span class="me-5">Show</span>
                            <ToggleSwitch
                                id="chartGridEnable"
                                :value="chartStyle.props.grid.vertical.show"
                                @input="setVerticalGridShow"
                                :width="50"
                                :height="25"
                            ></ToggleSwitch>
                          </div>
                          <div class="accordion-input d-flex justify-content-between">
                            <span>Style</span>
                            <select class="form-select form-select-sm" aria-label=".form-select-sm">
                              <option value="dash">Dash</option>
                              <option value="solid">Solid</option>
                            </select>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Candle Styling -->
              <div class="accordion-item">
                <h2 class="accordion-header" id="flush-headingTwo">
                  <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseTwo" aria-expanded="false" aria-controls="flush-collapseTwo">
                    Candle
                  </button>
                </h2>
                <div id="flush-collapseTwo" class="accordion-collapse collapse" aria-labelledby="flush-headingTwo" data-bs-parent="#accordionFlushExample">
                  <div class="accordion-body">
                    Placeholder content for this accordion, which is intended to demonstrate the
                    <code>.accordion-flush</code> class. This is the second item's accordion body.
                    Let's imagine this being filled with some actual content.
                  </div>
                </div>
              </div>

              <!-- Indicator Styling -->
              <div class="accordion-item">
                <h2 class="accordion-header" id="flush-headingThree">
                  <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseThree" aria-expanded="false" aria-controls="flush-collapseThree">
                    Indicators
                  </button>
                </h2>
                <div id="flush-collapseThree" class="accordion-collapse collapse" aria-labelledby="flush-headingThree" data-bs-parent="#accordionFlushExample">
                  <div class="accordion-body">Placeholder content for this accordion, which is intended to demonstrate the <code>.accordion-flush</code> class. This is the third item's accordion body. Nothing more exciting happening here in terms of content, but just filling up the space to make it look, at least at first glance, a bit more representative of how this would look in a real-world application.</div>
                </div>
              </div>

              <!-- X Axis Styling -->
              <div class="accordion-item">
                <h2 class="accordion-header" id="flush-headingFour">
                  <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseFour" aria-expanded="false" aria-controls="flush-collapseFour">
                    X Axis
                  </button>
                </h2>
                <div id="flush-collapseFour" class="accordion-collapse collapse" aria-labelledby="flush-headingFour" data-bs-parent="#accordionFlushExample">
                  <div class="accordion-body">Placeholder content for this accordion, which is intended to demonstrate the <code>.accordion-flush</code> class. This is the fourth item's accordion body. Nothing more exciting happening here in terms of content, but just filling up the space to make it look, at least at first glance, a bit more representative of how this would look in a real-world application.</div>
                </div>
              </div>

              <!-- Y Axis Styling -->
              <div class="accordion-item">
                <h2 class="accordion-header" id="flush-headingFive">
                  <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseFive" aria-expanded="false" aria-controls="flush-collapseFive">
                    Y Axis
                  </button>
                </h2>
                <div id="flush-collapseFive" class="accordion-collapse collapse" aria-labelledby="flush-headingFive" data-bs-parent="#accordionFlushExample">
                  <div class="accordion-body">Placeholder content for this accordion, which is intended to demonstrate the <code>.accordion-flush</code> class. This is the fifth item's accordion body. Nothing more exciting happening here in terms of content, but just filling up the space to make it look, at least at first glance, a bit more representative of how this would look in a real-world application.</div>
                </div>
              </div>

              <!-- Crosshair Styling -->
              <div class="accordion-item">
                <h2 class="accordion-header" id="flush-headingSix">
                  <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseSix" aria-expanded="false" aria-controls="flush-collapseSix">
                    Crosshair
                  </button>
                </h2>
                <div id="flush-collapseSix" class="accordion-collapse collapse" aria-labelledby="flush-headingSix" data-bs-parent="#accordionFlushExample">
                  <div class="accordion-body">Placeholder content for this accordion, which is intended to demonstrate the <code>.accordion-flush</code> class. This is the sixth item's accordion body. Nothing more exciting happening here in terms of content, but just filling up the space to make it look, at least at first glance, a bit more representative of how this would look in a real-world application.</div>
                </div>
              </div>

              <!-- Separator Styling -->
              <div class="accordion-item">
                <h2 class="accordion-header" id="flush-headingSeven">
                  <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseSeven" aria-expanded="false" aria-controls="flush-collapseSeven">
                    Separator
                  </button>
                </h2>
                <div id="flush-collapseSeven" class="accordion-collapse collapse" aria-labelledby="flush-headingThree" data-bs-parent="#accordionFlushExample">
                  <div class="accordion-body">Placeholder content for this accordion, which is intended to demonstrate the <code>.accordion-flush</code> class. This is the seventh item's accordion body. Nothing more exciting happening here in terms of content, but just filling up the space to make it look, at least at first glance, a bit more representative of how this would look in a real-world application.</div>
                </div>
              </div>

            </div>

          </div>

        </div>

        <!-- Screenshot Button -->
        <div class="icon-wrapper">
          <a type="button" class="d-flex align-items-center">
            <svg viewBox="0 0 1024 1024">
              <path d="M989.884071 853.441363H853.441745v136.442326a34.116311 34.116311 0 0 1-68.221163 0V853.441363h-579.874156a34.116311 34.116311 0 0 1-34.116311-34.116312v-579.874155H34.741949a34.116311 34.116311 0 1 1 0-68.221163h136.442326V34.741567a34.116311 34.116311 0 1 1 68.221163 0v136.442326h579.874155a34.116311 34.116311 0 0 1 34.116312 34.116311v579.874156h136.488166a34.116311 34.116311 0 0 1 0 68.221163z m-204.674949-614.024847H239.416898v545.792224h545.792224z m0 0"></path>
              <path d="M684.338523 614.386523q-52.131374-102.314555-151.271519-51.501076-27.813332-33.852732-66.570928-100.286141t-71.235134-153.506212q3.965147 83.119116 30.861682 136.373566c26.873615 53.208611 120.490057 209.717333 160.16445 240.727995s128.294292 16.043949 98.051449-71.796672z m-96.538733 19.871577c-34.43719-41.255869-33.715213-54.091028-13.190418-61.287885s41.656967-0.217739 65.230112 23.848184 28.558229 52.349113 13.511297 65.390552-31.136721 13.316478-65.573911-27.950851z m0 0"></path>
              <path d="M513.092287 537.157829l19.02354 26.243316c-33.669373 53.50657-70.432936 106.165102-91.679709 122.77059-39.662934 31.010661-128.294292 16.043949-98.051448-71.796672q52.165754-102.280175 151.271519-51.478156 9.110671-11.093245 19.413178-25.727618z m20.226836-30.254304q12.72056-19.974716 26.919454-44.304219 38.769057-66.467789 71.223674-153.494752-3.965147 83.119116-30.838762 136.373566c-9.649289 19.080839-27.870631 51.455236-49.117404 86.545645z m-94.40718 127.354575c34.46011-41.255869 33.738133-54.091028 13.213338-61.287885s-41.679887-0.217739-65.253032 23.848184-28.535309 52.349113-13.511297 65.390552 31.136721 13.316478 65.573911-27.950851z m0 0"></path>
            </svg>
            <span>Screenshot</span>
          </a>
        </div>

      </div>

      <!-- Chart Content -->
      <div class="chart-content">

        <!-- Drawing Tools -->
        <div class="chart-tools-bar">
          <div class="icon-wrapper" v-on:click="setShapeType('horizontalStraightLine')">
            <svg viewBox="0 0 24 24">
              <rect x="5" y="11.5" width="14" height="1" rx="0.5"></rect>
              <ellipse cx="12" cy="12" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('horizontalRayLine')">
            <svg viewBox="0 0 24 24">
              <rect x="4.5" y="11.5" width="15" height="1" rx="0.5"></rect>
              <ellipse cx="6" cy="12" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="13" cy="12" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('horizontalSegment')">
            <svg viewBox="0 0 24 24">
              <rect x="5" y="11.5" width="14" height="1" rx="0.5"></rect>
              <ellipse cx="6" cy="12" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="18" cy="12" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('verticalStraightLine')">
            <svg viewBox="0 0 24 24">
              <rect x="11.5" y="4" width="1" height="16" rx="0.5"></rect>
              <ellipse cx="12" cy="12" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('verticalRayLine')">
            <svg viewBox="0 0 24 24">
              <rect x="11.5" y="4.5" width="1" height="15" rx="0.5"></rect>
              <ellipse cx="12" cy="18" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="12" cy="11" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('verticalSegment')">
            <svg viewBox="0 0 24 24">
              <rect x="11.5" y="5" width="1" height="14" rx="0.5"></rect>
              <ellipse cx="12" cy="18" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="12" cy="6" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('straightLine')">
            <svg viewBox="0 0 24 24">
              <rect x="5.989593505859375" y="17.303306579589844" width="16" height="1" rx="0.5"
                    transform="matrix(0.7071067690849304,-0.7071067690849304,0.7071067690849304,0.7071067690849304,-10.480973816180722,9.303303481670355)"
              ></rect>
              <ellipse cx="14" cy="10" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="10" cy="14" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('rayLine')">
            <svg viewBox="0 0 24 24">
              <rect x="6.989593505859375" y="16.303314208984375" width="15" height="1" rx="0.5"
                    transform="matrix(0.7071067690849304,-0.7071067690849304,0.7071067690849304,0.7071067690849304,-9.480979210977239,9.71751925443823)"
              ></rect>
              <ellipse cx="13" cy="11" rx="1.5" ry="1.5"></ellipse><ellipse cx="8" cy="16" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('segment')">
            <svg viewBox="0 0 24 24"><rect x="5.989593505859375" y="17.303298950195312" width="14" height="1" rx="0.5"
                                           transform="matrix(0.7071067690849304,-0.7071067690849304,0.7071067690849304,0.7071067690849304,-10.480968421384205,9.30330124707234)"
            ></rect>
              <ellipse cx="16" cy="8" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="7" cy="17" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('priceLine')">
            <svg viewBox="0 0 24 24">
              <rect x="4.5" y="13.5" width="15" height="1" rx="0.5"></rect>
              <ellipse cx="6" cy="14" rx="1.5" ry="1.5"></ellipse>
              <path d="M8.314036947998046,12.235949340820312L10.985912947998047,12.235949340820312L10.985912947998047,11.517199340820312L10.151922947998047,11.517199340820312L10.151922947998047,7.735949340820312L9.497632947998047,7.735949340820312C9.214422947998047,7.917589340820312,8.915602947998046,8.030869340820313,8.464427947998047,8.108999340820313L8.464427947998047,8.661729340820312L9.274972947998046,8.661729340820312L9.274972947998046,11.517199340820312L8.314036947998046,11.517199340820312L8.314036947998046,12.235949340820312ZM11.581612947998046,12.235949340820312L14.556222947998048,12.235949340820312L14.556222947998048,11.493759340820311L13.597242947998048,11.493759340820311C13.386302947998047,11.493759340820311,13.093332947998046,11.517199340820312,12.864822947998046,11.546499340820311C13.675362947998046,10.724229340820312,14.347242947998048,9.831649340820313,14.347242947998048,9.001579340820312C14.347242947998048,8.151969340820312,13.788642947998046,7.610949340820312,12.948802947998047,7.610949340820312C12.343332947998046,7.610949340820312,11.946852947998046,7.845329340820312,11.532782947998047,8.290639340820313L12.024972947998048,8.778919340820313C12.247632947998046,8.525009340820313,12.511302947998047,8.308219340820312,12.835522947998047,8.308219340820312C13.261302947998047,8.308219340820312,13.501532947998047,8.593369340820313,13.501532947998047,9.044539340820313C13.501532947998047,9.757429340820313,12.792552947998047,10.618759340820311,11.581612947998046,11.726179340820313L11.581612947998046,12.235949340820312ZM16.460522947998047,12.360949340820312C17.312082947998046,12.360949340820312,18.026902947998046,11.894149340820313,18.026902947998046,11.048449340820312C18.026902947998046,10.431259340820311,17.642162947998045,10.050399340820313,17.144112947998046,9.911729340820312L17.144112947998046,9.882429340820313C17.612862947998046,9.696889340820313,17.882402947998045,9.331649340820313,17.882402947998045,8.823839340820312C17.882402947998045,8.032829340820312,17.300362947998046,7.610949340820312,16.44294294799805,7.610949340820312C15.921462947998046,7.610949340820312,15.495682947998047,7.821889340820313,15.110912947998047,8.151969340820312L15.565992947998048,8.722279340820313C15.825752947998048,8.460559340820312,16.083572947998046,8.308219340820312,16.401922947998045,8.308219340820312C16.77888294799805,8.308219340820312,16.99568294799805,8.525009340820313,16.99568294799805,8.892199340820312C16.99568294799805,9.319929340820313,16.730052947998047,9.610949340820312,15.921462947998046,9.610949340820312L15.921462947998046,10.247669340820313C16.88044294799805,10.247669340820313,17.138252947998048,10.530869340820313,17.138252947998048,10.991809340820312C17.138252947998048,11.407829340820314,16.833572947998046,11.642199340820312,16.38239294799805,11.642199340820312C15.974192947998047,11.642199340820312,15.657782947998047,11.433219340820312,15.392162947998047,11.161729340820312L14.978102947998046,11.743759340820311C15.290602947998046,12.097279340820313,15.765212947998048,12.360949340820312,16.460522947998047,12.360949340820312Z"></path>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('parallelStraightLine')">
            <svg viewBox="0 0 24 24">
              <rect x="7.989593505859375" y="20.303314208984375" width="16" height="1" rx="0.5"
                    transform="matrix(0.7071067690849304,-0.7071067690849304,0.7071067690849304,0.7071067690849304,-12.016513056401891,11.596198947183439)"
              ></rect>
              <rect x="3.4586830139160156" y="15.292892456054688" width="16" height="1" rx="0.5"
                    transform="matrix(0.7071067690849304,-0.7071067690849304,0.7071067690849304,0.7071067690849304,-9.800682931907204,6.924842852749634)"
              ></rect>
              <ellipse cx="16" cy="13" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="12" cy="17" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="9.5" cy="10" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('priceChannelLine')">
            <svg viewBox="0 0 24 24">
              <rect x="5.989593505859375" y="17.303298950195312" width="16" height="1" rx="0.5"
                    transform="matrix(0.7071067690849304,-0.7071067690849304,0.7071067690849304,0.7071067690849304,-10.480968421384205,9.30330124707234)"
              ></rect>
              <rect x="5.031974792480469" y="13.607009887695312" width="12" height="1" rx="0.5"
                    transform="matrix(0.7071067690849304,-0.7071067690849304,-0.7071067690849304,-0.7071067690849304,11.095440153447726,26.786762123917924)"
              ></rect>
              <rect x="11.40380859375" y="19.303298950195312" width="12" height="1" rx="0.5"
                    transform="matrix(0.7071067690849304,-0.7071067690849304,-0.7071067690849304,-0.7071067690849304,16.98959169711361,41.016502553537975)"
              ></rect>
              <ellipse cx="14" cy="10" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="10" cy="14" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="15" cy="15" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('fibonacciLine')">
            <svg viewBox="0 0 24 24">
              <rect x="4" y="6" width="16" height="1" rx="0.5"></rect>
              <rect x="4" y="9" width="16" height="1" rx="0.5"></rect>
              <rect x="4" y="15" width="16" height="1" rx="0.5"></rect>
              <rect x="4" y="18" width="16" height="1" rx="0.5"></rect>
              <rect x="4" y="12" width="16" height="1" rx="0.5"></rect>
              <ellipse cx="12" cy="18.5" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="16" cy="6.5" rx="1.5" ry="1.5"></ellipse>
              <ellipse cx="8" cy="6.5" rx="1.5" ry="1.5"></ellipse>
            </svg>
          </div>
          <hr class="divider"/>
          <div class="icon-wrapper" v-on:click="setShapeType('rect')">
            <svg width="24px" height="24px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path d="M5 3H3v2h2V3zm2 0h2v2H7V3zm6 0h-2v2h2V3zm2 0h2v2h-2V3zm4 0h2v2h-2V3zM3 7h2v2H3V7zm2 4H3v2h2v-2zm-2 4h2v2H3v-2zm2 4H3v2h2v-2zm2 0h2v2H7v-2zm6 0h-2v2h2v-2zm6-8h2v2h-2v-2zm2-4h-2v2h2V7zm-6 10v-2h6v2h-2v2h-2v2h-2v-4zm4 2v2h2v-2h-2z" />
            </svg>
          </div>
          <div class="icon-wrapper" v-on:click="setShapeType('circle')" style="padding: 2px;">
            <svg viewBox="0 0 59.94 59.94">
              <path d="M30.694,59.94c-1.078,0-1.967-0.857-1.998-1.941c-0.032-1.104,0.837-2.025,1.94-2.058c0.314-0.009,0.628-0.022,0.939-0.042
          c1.12-0.084,2.051,0.771,2.119,1.873c0.068,1.103-0.771,2.052-1.873,2.119c-0.354,0.021-0.711,0.037-1.068,0.048
          C30.733,59.94,30.714,59.94,30.694,59.94z M24.696,59.468c-0.121,0-0.244-0.011-0.368-0.034c-0.356-0.066-0.71-0.139-1.062-0.219
          c-1.077-0.245-1.752-1.317-1.507-2.394c0.245-1.077,1.319-1.751,2.394-1.507c0.301,0.068,0.604,0.131,0.907,0.188
          c1.086,0.202,1.803,1.246,1.6,2.332C26.481,58.796,25.64,59.468,24.696,59.468z M37.63,58.88c-0.872,0-1.673-0.574-1.923-1.455
          c-0.302-1.063,0.314-2.168,1.378-2.47c0.299-0.085,0.595-0.176,0.89-0.271c1.054-0.339,2.179,0.233,2.521,1.284
          c0.341,1.05-0.233,2.179-1.284,2.521c-0.342,0.111-0.687,0.216-1.034,0.314C37.995,58.855,37.81,58.88,37.63,58.88z M18.013,57.318
          c-0.284,0-0.572-0.061-0.847-0.188c-0.327-0.153-0.651-0.312-0.972-0.477c-0.982-0.506-1.369-1.712-0.863-2.693
          c0.505-0.981,1.711-1.368,2.693-0.863c0.275,0.142,0.555,0.278,0.837,0.41c1,0.468,1.432,1.659,0.964,2.659
          C19.486,56.892,18.765,57.318,18.013,57.318z M44.087,56.122c-0.688,0-1.356-0.354-1.729-0.991
          c-0.558-0.953-0.236-2.178,0.718-2.735c0.267-0.156,0.53-0.318,0.791-0.484c0.934-0.594,2.168-0.318,2.762,0.613
          c0.593,0.932,0.318,2.168-0.613,2.762c-0.304,0.193-0.61,0.381-0.922,0.563C44.777,56.034,44.429,56.122,44.087,56.122z
          M12.08,53.564c-0.449,0-0.9-0.15-1.273-0.458c-0.277-0.229-0.551-0.464-0.819-0.703c-0.825-0.734-0.898-1.998-0.163-2.823
          c0.732-0.824,1.998-0.899,2.823-0.163c0.231,0.206,0.468,0.407,0.708,0.605c0.852,0.704,0.971,1.965,0.268,2.816
          C13.227,53.317,12.655,53.564,12.08,53.564z M49.644,51.843c-0.514,0-1.026-0.196-1.417-0.589c-0.779-0.782-0.777-2.049,0.006-2.828
          c0.222-0.221,0.439-0.445,0.654-0.674c0.757-0.804,2.021-0.843,2.827-0.087c0.805,0.756,0.844,2.022,0.087,2.827
          c-0.244,0.26-0.493,0.516-0.746,0.768C50.666,51.649,50.155,51.843,49.644,51.843z M7.271,48.456c-0.617,0-1.227-0.284-1.618-0.821
          c-0.211-0.29-0.418-0.585-0.618-0.881c-0.617-0.916-0.376-2.159,0.539-2.777c0.917-0.618,2.159-0.376,2.777,0.539
          c0.173,0.257,0.351,0.511,0.534,0.762c0.65,0.893,0.455,2.144-0.438,2.795C8.093,48.331,7.679,48.456,7.271,48.456z M53.967,46.316
          c-0.349,0-0.701-0.091-1.022-0.282c-0.949-0.566-1.259-1.794-0.693-2.742c0.16-0.269,0.315-0.539,0.466-0.813
          c0.531-0.968,1.747-1.322,2.717-0.789c0.968,0.532,1.321,1.748,0.789,2.717c-0.174,0.315-0.353,0.627-0.536,0.935
          C55.311,45.968,54.648,46.316,53.967,46.316z M3.882,42.312c-0.797,0-1.549-0.479-1.86-1.265c-0.132-0.332-0.259-0.669-0.379-1.008
          c-0.369-1.041,0.175-2.185,1.216-2.554c1.039-0.369,2.184,0.174,2.554,1.216c0.104,0.294,0.214,0.584,0.328,0.873
          c0.407,1.026-0.096,2.188-1.123,2.596C4.376,42.266,4.127,42.312,3.882,42.312z M56.779,39.885c-0.186,0-0.374-0.026-0.562-0.081
          c-1.06-0.31-1.669-1.42-1.359-2.48c0.087-0.299,0.17-0.601,0.248-0.904c0.272-1.068,1.357-1.721,2.433-1.442
          c1.069,0.272,1.716,1.362,1.442,2.433c-0.088,0.347-0.184,0.691-0.283,1.035C58.443,39.318,57.645,39.885,56.779,39.885z
          M2.133,35.515c-0.992,0-1.853-0.738-1.981-1.748c-0.045-0.357-0.084-0.716-0.115-1.078c-0.097-1.101,0.718-2.07,1.818-2.166
          c1.122-0.112,2.07,0.719,2.166,1.818c0.027,0.31,0.061,0.617,0.1,0.922c0.139,1.096-0.637,2.097-1.732,2.236
          C2.302,35.51,2.217,35.515,2.133,35.515z M57.895,32.956c-0.023,0-0.047,0-0.071-0.001c-1.104-0.039-1.967-0.965-1.929-2.069
          c0.011-0.311,0.017-0.622,0.017-0.935v-0.113c0-1.104,0.896-2,2-2s2,0.896,2,2v0.113c0,0.359-0.006,0.718-0.019,1.075
          C59.855,32.106,58.968,32.956,57.895,32.956z M2.125,28.497c-0.082,0-0.164-0.005-0.247-0.015c-1.096-0.135-1.876-1.133-1.741-2.229
          c0.044-0.356,0.094-0.711,0.149-1.062c0.172-1.092,1.198-1.837,2.286-1.665c1.092,0.172,1.837,1.195,1.665,2.286
          c-0.048,0.308-0.092,0.617-0.13,0.929C3.982,27.754,3.12,28.497,2.125,28.497z M57.249,25.855c-0.918,0-1.745-0.636-1.951-1.57
          c-0.067-0.303-0.14-0.605-0.217-0.904c-0.278-1.069,0.363-2.161,1.433-2.438c1.072-0.28,2.161,0.364,2.438,1.433
          c0.091,0.348,0.175,0.697,0.252,1.051c0.237,1.078-0.444,2.146-1.523,2.383C57.536,25.84,57.391,25.855,57.249,25.855z M3.84,21.693
          c-0.242,0-0.49-0.045-0.729-0.14c-1.028-0.403-1.535-1.563-1.131-2.592c0.132-0.337,0.27-0.671,0.413-1
          c0.441-1.012,1.616-1.476,2.633-1.033c1.013,0.441,1.475,1.62,1.033,2.633c-0.124,0.284-0.242,0.571-0.356,0.861
          C5.393,21.211,4.638,21.693,3.84,21.693z M54.887,19.249c-0.732,0-1.438-0.402-1.788-1.101c-0.139-0.275-0.283-0.546-0.433-0.814
          c-0.536-0.966-0.188-2.184,0.778-2.72c0.966-0.534,2.183-0.187,2.72,0.778c0.175,0.315,0.345,0.635,0.507,0.957
          c0.497,0.987,0.1,2.189-0.887,2.686C55.496,19.18,55.189,19.249,54.887,19.249z M7.206,15.532c-0.406,0-0.815-0.123-1.17-0.379
          C5.14,14.505,4.94,13.255,5.587,12.36c0.211-0.292,0.428-0.58,0.65-0.865c0.682-0.871,1.938-1.022,2.808-0.343
          c0.87,0.681,1.023,1.938,0.343,2.808c-0.191,0.245-0.377,0.493-0.559,0.744C8.437,15.245,7.826,15.532,7.206,15.532z M50.95,13.445
          c-0.557,0-1.109-0.23-1.505-0.682c-0.205-0.233-0.413-0.461-0.625-0.686c-0.761-0.802-0.727-2.067,0.075-2.827
          s2.068-0.727,2.827,0.075c0.249,0.262,0.492,0.528,0.73,0.801c0.729,0.831,0.645,2.095-0.186,2.822
          C51.887,13.282,51.418,13.445,50.95,13.445z M12.002,10.404c-0.574,0-1.145-0.246-1.54-0.723C9.757,8.831,9.874,7.57,10.723,6.865
          c0.275-0.229,0.555-0.451,0.837-0.67c0.871-0.676,2.127-0.518,2.806,0.356c0.677,0.873,0.518,2.129-0.356,2.806
          c-0.247,0.191-0.491,0.387-0.731,0.586C12.904,10.253,12.452,10.404,12.002,10.404z M45.689,8.796c-0.389,0-0.781-0.113-1.126-0.349
          c-0.257-0.175-0.516-0.345-0.778-0.51c-0.936-0.588-1.217-1.823-0.629-2.758c0.588-0.936,1.825-1.217,2.758-0.629
          c0.306,0.192,0.607,0.39,0.905,0.594c0.912,0.623,1.146,1.867,0.523,2.779C46.956,8.492,46.328,8.796,45.689,8.796z M17.915,6.631
          c-0.749,0-1.468-0.423-1.81-1.146c-0.472-0.999-0.045-2.191,0.954-2.663c0.323-0.152,0.649-0.3,0.978-0.44
          c1.015-0.44,2.191,0.032,2.627,1.047c0.437,1.015-0.032,2.191-1.047,2.627c-0.285,0.123-0.568,0.251-0.849,0.384
          C18.492,6.57,18.201,6.631,17.915,6.631z M39.439,5.605c-0.225,0-0.453-0.038-0.676-0.119c-0.293-0.104-0.589-0.205-0.888-0.301
          c-1.052-0.336-1.633-1.461-1.297-2.514c0.336-1.052,1.459-1.634,2.514-1.297c0.344,0.109,0.685,0.226,1.022,0.348
          c1.04,0.373,1.58,1.519,1.206,2.558C41.028,5.096,40.26,5.605,39.439,5.605z M24.58,4.454c-0.941,0-1.78-0.667-1.963-1.626
          c-0.206-1.085,0.506-2.132,1.591-2.339c0.352-0.066,0.706-0.129,1.062-0.183c1.093-0.164,2.112,0.582,2.279,1.674
          s-0.582,2.112-1.674,2.279c-0.309,0.048-0.614,0.101-0.919,0.159C24.83,4.442,24.704,4.454,24.58,4.454z M32.59,4.076
          c-0.063,0-0.126-0.003-0.189-0.009c-0.309-0.029-0.618-0.053-0.931-0.07c-1.103-0.064-1.944-1.011-1.881-2.113
          c0.064-1.103,1.004-1.936,2.113-1.881c0.359,0.021,0.718,0.049,1.073,0.082c1.1,0.104,1.907,1.079,1.804,2.179
          C34.481,3.299,33.61,4.076,32.59,4.076z"
              />
            </svg>
          </div>
          <hr class="divider"/>
          <div class="icon-wrapper" v-on:click="removeAllShape()">
            <svg viewBox="0 0 1024 1024" style="width: 34px; height: 34px;">
              <path d="M256 333.872a28.8 28.8 0 0 1 28.8 28.8V768a56.528 56.528 0 0 0 56.544 56.528h341.328A56.528 56.528 0 0 0 739.2 768V362.672a28.8 28.8 0 0 1 57.6 0V768a114.128 114.128 0 0 1-114.128 114.128H341.328A114.128 114.128 0 0 1 227.2 768V362.672a28.8 28.8 0 0 1 28.8-28.8zM405.344 269.648a28.8 28.8 0 0 0 28.8-28.8 56.528 56.528 0 0 1 56.528-56.544h42.656a56.528 56.528 0 0 1 56.544 56.544 28.8 28.8 0 0 0 57.6 0 114.128 114.128 0 0 0-112.64-114.128h-45.648a114.144 114.144 0 0 0-112.64 114.128 28.8 28.8 0 0 0 28.8 28.8z"></path>
              <path d="M163.2 266.672a28.8 28.8 0 0 1 28.8-28.8h640a28.8 28.8 0 0 1 0 57.6H192a28.8 28.8 0 0 1-28.8-28.8zM426.672 371.2a28.8 28.8 0 0 1 28.8 28.8v320a28.8 28.8 0 0 1-57.6 0V400a28.8 28.8 0 0 1 28.8-28.8zM597.344 371.2a28.8 28.8 0 0 1 28.8 28.8v320a28.8 28.8 0 0 1-57.6 0V400a28.8 28.8 0 0 1 28.8-28.8z"></path>
            </svg>
          </div>
        </div>

        <!-- Chart -->
        <div id="chart" class="chart" v-bind:style="{backgroundColor: themeStore.theme === 'dark' ? chartStyle.themes.dark.chartBackgroundColor : chartStyle.themes.light.chartBackgroundColor}"/>

      </div>

    </div>

    <!-- Orderbook -->
    <Orderbook
        v-bind:midpoint="this.midpoint"
        v-bind:asks="this.asks"
        v-bind:bids="this.bids"
        v-bind:book-initialized="this.bookInitialized"
    ></Orderbook>


</div>
</template>
<script>
import {dispose, init} from "klinecharts";
import {themeStore} from "../store/themeStore";
import ToggleSwitch from "../components/ToggleSwitch";
import ColorPickerInput from "../components/ColorPickerInput";
import {chartStyle} from "../store/chartStyleStore";
import {emojiTechnicalIndicator} from "../components/klinechart/indicators/emojiIndicator.js"
import { checkCoordinateOnSegment } from "klinecharts/lib/shape/shapeHelper"
import StatusIndicator from "../components/StatusIndicator";
import jwtInterceptor from "../shared/jwt.interceptor";
import Dropdown from "../components/Dropdown";
import {websocket} from "../websocket/websocket";
import Orderbook from "../components/orderbook/Orderbook";

const rect = {
  name: 'rect',
  totalStep: 3,
  checkEventCoordinateOnShape: ({ dataSource, eventCoordinate }) => {
    return checkCoordinateOnSegment(dataSource[0], dataSource[1], eventCoordinate)
  },
  createShapeDataSource: ({ coordinates }) => {
    if (coordinates.length === 2) {
      return [
        {
          type: 'line',
          isDraw: false,
          isCheck: true,
          dataSource: [
            [{ ...coordinates[0] }, { x: coordinates[1].x, y: coordinates[0].y }],
            [{ x: coordinates[1].x, y: coordinates[0].y }, { ...coordinates[1] }],
            [{ ...coordinates[1] }, { x: coordinates[0].x, y: coordinates[1].y }],
            [{ x: coordinates[0].x, y: coordinates[1].y }, { ...coordinates[0] }]
          ]
        },
        {
          type: 'polygon',
          isDraw: true,
          isCheck: false,
          styles: { style: 'fill' },
          dataSource: [[
            { ...coordinates[0] },
            { x: coordinates[1].x, y: coordinates[0].y },
            { ...coordinates[1] },
            { x: coordinates[0].x, y: coordinates[1].y }
          ]]
        },
        {
          type: 'polygon',
          isDraw: true,
          isCheck: false,
          dataSource: [[
            { ...coordinates[0] },
            { x: coordinates[1].x, y: coordinates[0].y },
            { ...coordinates[1] },
            { x: coordinates[0].x, y: coordinates[1].y }
          ]]
        }
      ]
    }
    return []
  }
}
const circle = {
  name: 'circle',
  totalStep: 3,
  checkEventCoordinateOnShape: ({ dataSource, eventCoordinate }) => {
    const xDis = Math.abs(dataSource.x - eventCoordinate.x)
    const yDis = Math.abs(dataSource.y - eventCoordinate.y)
    const r = Math.sqrt(xDis * xDis + yDis * yDis)
    return Math.abs(r - dataSource.radius) < 3;
  },
  createShapeDataSource: ({ coordinates }) => {
    if (coordinates.length === 2) {
      const xDis = Math.abs(coordinates[0].x - coordinates[1].x)
      const yDis = Math.abs(coordinates[0].y - coordinates[1].y)
      const radius = Math.sqrt(xDis * xDis + yDis * yDis)
      return [
        {
          type: 'arc',
          isDraw: true,
          isCheck: false,

          dataSource: [
            { ...coordinates[0], radius, startAngle: 0, endAngle: Math.PI * 2 }
          ]
        },
        {
          type: 'arc',
          isDraw: true,
          isCheck: true,
          dataSource: [
            { ...coordinates[0], radius, startAngle: 0, endAngle: Math.PI * 2 }
          ]
        }
      ];
    }
    return []
  }
}


export default {
  name: 'Chart',
  components: {
    Orderbook,
    Dropdown,
    StatusIndicator,
    ColorPickerInput,
    ToggleSwitch
  },
  data: function () {
    return {
      websocket,
      themeStore,
      chartStyle,
      mainTechnicalIndicatorTypes: [
        {
          name: "MA",
          checked: false,
        },
        {
          name: "EMA",
          checked: false,
        },
        {
          name: "SAR",
          checked: false,
        },

      ],
      subTechnicalIndicatorTypes: [
        {
          name: "VOL",
          checked: false,
        },
        {
          name: "MACD",
          checked: true,
        },
        {
          name: "KDJ",
          checked: true,
        },
        {
          name: "EMOJI",
          checked: true,
        },
      ],
      customEmojiIndicatorTypes: {
        name: 'EMOJI',
        checked: false
      },
      shapes: [
        /* Default Shapes */
        { key: 'horizontalRayLine', text: 'custom circle' },
        { key: 'horizontalSegment', text: 'custom circle' },
        { key: 'horizontalStraightLine', text: 'custom circle' },
        { key: 'verticalRayLine', text: 'custom circle' },
        { key: 'verticalSegment', text: 'custom circle' },
        { key: 'verticalStraightLine', text: 'custom circle' },
        { key: 'rayLine', text: 'custom circle' },
        { key: 'segment', text: 'custom circle' },
        { key: 'straightLine', text: 'custom circle' },
        { key: 'priceLine', text: 'price line' },
        { key: 'priceChannelLine', text: 'price channel line' },
        { key: 'parallelStraightLine', text: 'parallel lines' },
        { key: 'fibonacciLine', text: 'fibonacci retracement' },

        /* Custom Shapes */
        { key: 'rect', text: 'custom rectangle' },
        { key: 'circle', text: 'custom circle' }
      ],

      /* Chart & Websocket data */
      dataConnected: false,
      simTimer: undefined,
      simTimout: 1000,
      chartPeriodTimerFunc: undefined,
      chartPeriod: 3600, /*1hr*/
      heartbeat: {},

      /* Orderbook data */
      asks: [],
      bids: [],
      midpoint: "",
      spread: "",
      bookInitialized: false,
      bookArrayLimit: 100, // size of each bid/ask array

      /* Exchange Select */
      selectedExchange: '',
      exchanges: [],
      selectedProduct: '',
      productInputDisabled: true,
      products: [],

      /* Period Select */
      templates: ['1m','5m','30m','1h', '6h', '1d', '1w'],
      category: '5m'


    }
  },
  watch: {
    themeStore: {
      handler(newValue, oldValue) {
        // this will be run immediately on component creation.
        this.setTheme(newValue.theme)
      },
      deep: true, // fire on all nested mutations
    },

  },
  computed: {
    date() {
      let date = new Date(this.heartbeat.time)
      let year = date.getFullYear();
      let month = date.getMonth()+1;
      let dt = date.getDate();
      if (dt < 10) {
        dt = '0' + dt;
      }
      if (month < 10) {
        month = '0' + month;
      }
      return (year+'-' + month + '-'+dt)
    },
    time() {
      let date = new Date(this.heartbeat.time);
      return date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds() + ':' + date.getMilliseconds()
    },

  },
  created() {
    chartStyle.getStyle()
    window.addEventListener("resize", this.handleWindowResize);
  },
  mounted: function () {
    this.getExchanges()
    this.kLineChart = init('chart')
    this.setTheme(themeStore.theme)
    this.kLineChart.addShapeTemplate([rect, circle])
    this.kLineChart.addTechnicalIndicatorTemplate(emojiTechnicalIndicator)
    this.volPaneId = this.kLineChart.createTechnicalIndicator('VOL', false)
    this.macdPaneId = this.kLineChart.createTechnicalIndicator('MACD', false)
    this.kdjPaneId = this.kLineChart.createTechnicalIndicator('KDJ', false)
    this.emojiPaneId = this.kLineChart.createTechnicalIndicator('EMOJI', false)
    this.subTechnicalIndicatorTypes.forEach((item) => {
      this.setSubTechnicalIndicator(item)
    })
    this.kLineChart.setPriceVolumePrecision(4, 10)
    this.$bus.on("onWebsocketMessage", data => {
      // hack for api sending multiple json objects in one message
      let message;
      let messages;
      let isMultiple = false
      try {
        message = JSON.parse(data)
      } catch(e) {
        try {
          messages = data.split(/(?!})(?={)/).map(function(v, i) { return JSON.parse(v); });
          isMultiple = true
        } catch (e) {
          console.log("error parsing message: " + e.toString())
        }
      }
      if (isMultiple) {
        messages.forEach(msg => {
          this.handleMessage(msg)
        });
      } else {
        this.handleMessage(message)
      }
    })
    this.$bus.on("onWebsocketError", msg => {
      console.log(msg)
    })

  },
  beforeRouteUpdate(to, from, next) {
    this.stopChartTimer()
    this.socketSendUnsubReq()
    next()
  },
  methods: {

    /* Chart functions */
    setChartType: function (type) {
      this.kLineChart.setStyleOptions({
        candle: {
          type
        }
      })
    },
    setTheme: function (theme) {
      console.log("setting theme: " + theme)
      if (theme === "dark" || theme === "light") {
        chartStyle.setThemeOptions(theme)
        this.saveStyleOptions()
      }
    },
    saveStyleOptions: function () {
      this.kLineChart.setStyleOptions(chartStyle.props)
      chartStyle.saveStyle()
    },
    setCandleTechnicalIndicator: function (type) {
      type.checked = !type.checked // starts out as default false, then passed thru need to flip here
      if (type.checked) {
        this.kLineChart.createTechnicalIndicator(type.name, true, { id: 'candle_pane' })
      } else {
        this.kLineChart.removeTechnicalIndicator('candle_pane', type.name)
      }
    },
    setSubTechnicalIndicator: function (type) {
      type.checked = !type.checked // starts out as default false, then passed thru need to flip here
      switch (type.name) {
        case "VOL":
          if (type.checked) {
            this.kLineChart.createTechnicalIndicator(type.name, true, { id: this.volPaneId })
          } else {
            this.kLineChart.removeTechnicalIndicator(this.volPaneId, type.name)
          }
          break;
        case "MACD":
          if (type.checked) {
            this.kLineChart.createTechnicalIndicator(type.name, true, { id: this.macdPaneId })
          } else {
            this.kLineChart.removeTechnicalIndicator(this.macdPaneId, type.name)
          }
          break;
        case "KDJ":
          if (type.checked) {
            this.kLineChart.createTechnicalIndicator(type.name, true, { id: this.kdjPaneId })
          } else {
            this.kLineChart.removeTechnicalIndicator(this.kdjPaneId, type.name)
          }
          break;
        case "EMOJI":
          if (type.checked) {
            this.kLineChart.createTechnicalIndicator(type.name, true, { id: this.emojiPaneId })
          } else {
            this.kLineChart.removeTechnicalIndicator(this.emojiPaneId, type.name)
          }
          break;
      }
    },
    setShapeType: function (shapeName) {
      this.kLineChart.createShape(shapeName)
    },
    removeAllShape () {
      console.log("removing all shapes")
      this.kLineChart.removeShape()
    },
    setGridShow: function (enabled) {
      console.log("setting grid show: " + enabled)
      chartStyle.props.grid.show = enabled
      this.saveStyleOptions()
    },
    setHorizontalGridShow: function (enabled) {
      console.log("setting horizontal grid show: " + enabled)
      chartStyle.props.grid.horizontal.show = enabled
      this.saveStyleOptions()
    },
    setGridHorizontalStyle: function (style) {
      console.log("setting horizontal grid style: " + style)
      let index = style.target.options.selectedIndex
      chartStyle.props.grid.horizontal.style = style.target.options[index].value
      this.saveStyleOptions()
    },
    setGridHorizontalSize: function (size) {
      console.log("setting horizontal grid size: " + size)
      chartStyle.props.grid.horizontal.size = size.target.value
      this.saveStyleOptions()
    },
    setGridHorizontalColor: function (color) {
      console.log("setting horizontal grid color: " + color)
      chartStyle.props.grid.horizontal.color = color
      this.saveStyleOptions()
    },
    setVerticalGridShow: function (enabled) {
      console.log("setting vertical grid show: " + enabled)
      chartStyle.props.grid.horizontal.show = enabled
      this.saveStyleOptions()
    },
    handleWindowResize(e) {
      this.kLineChart.resize()
    },

    /* Exchange and Product functions */
    onExchangeSelected(keyword) {
      this.selectedExchange = keyword;
      this.getProducts(keyword.id)
      console.log(keyword.id + " has been selected");
    },
    getExchangeValues(keyword) {
      console.log("You could refresh options by querying the API with " + keyword);
    },
    onProductSelected(selection) {
      this.asks = []
      this.bids = []
      this.selectedProduct = selection;
      this.socketSendSubReq()
      this.dataConnected = true


    },
    onPeriodChanged(selection) {
      this.getCandles()
    },
    getProductValues(keyword) {
      //console.log("You could refresh options by querying the API with " + keyword);
    },

    /* REST endpoints */
    async getExchanges() {
      const response = await jwtInterceptor.get("/api/exchanges", {
        withCredentials: true,
        credentials: "include",
        headers: {
          'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
        },
      });
      if (response && response.data) {
        this.handleMessage(response.data)
      }
    },
    async getCandles() {
      const response = await jwtInterceptor.get("/api/exchanges/" + this.selectedExchange.id + "/products/" + this.selectedProduct.id + "/candles", {
        params: {
          granularity: this.category
        },
        withCredentials: true,
        credentials: "include",
        headers: {
          'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
        },
      });
      if (response && response.data) {
        //this.startChartTimer()
        this.handleMessage(response.data)
      }
    },
    async getProducts(selection) {
      const response = await jwtInterceptor.get("/api/exchanges/" + selection + "/products", {
        withCredentials: true,
        credentials: "include",
        headers: {
          'X-Requested-With': 'XMLHttpRequest' // CSRF prevention
        },
      });
      if (response && response.data) {
        this.handleMessage(response.data)
      }
    },

    /* Data stream functions */
    socketSendSubReq() {
      let subRequest = {
        "action": "subscribe",
        "exchange_id": this.selectedExchange.id,
        "message": this.selectedProduct.id
      }
      this.websocket.socket.send(JSON.stringify(subRequest));
      this.getCandles()
    },
    socketSendUnsubReq() {
      let unsubRequest = {
        "action": "unsubscribe",
        "exchange_id": this.selectedExchange.id,
        "message": this.selectedProduct.id
      }
      this.websocket.socket.send(JSON.stringify(unsubRequest));
    },
    handleMessage(message) {
      switch (message.type) {
        case "welcome":
          console.log(message.message)
          break;
        case "error":
          console.log(message.message)
          break;
        case "l2update":
          this.handleL2Update(message)
          break;
        case "heartbeat":
          this.handleHeartbeat(message)
          break;
        case "match":
          this.handleMatch(message)
          break;
        case "snapshot":
          this.handleSnapshot(message)
          break;
        case "getProductBook":
          this.handleGetProductBook(message)
          break;
        case "getProductCandles":
          this.handleGetProductCandles(message)
          break;
        case "getProducts":
          this.handleGetProducts(message)
          break;
        case "getExchanges":
          this.handleGetExchanges(message)
          break;
      }
    },
    handleL2Update(message) {
      /*
      Response object format
      {
        "type": "l2update",
        "product_id": "BTC-USD",
        "time": "2022-04-29T14:00:48.180681Z",
        "changes": [
            [
                "sell",
                "38974.72",
                "0.00000000"
            ]
        ]
      }*/
      if (!this.bookInitialized) {
        return;
      }
      message.changes.forEach(change => {
        this.update_book(change, change[0]);
      });
    },
    handleHeartbeat(message) {
      /*
      Response object format
      {
        "type": "heartbeat",
        "sequence": 541837244,
        "last_trade_id": 39067179,
        "product_id": "BTC-USD",
        "time": "2022-04-29T14:00:48.994274Z"
      }*/
      this.heartbeat = message
    },
    handleMatch(message) {
      /*
      Response object format
      {
        "type": "match",
        "trade_id": 10,
        "sequence": 50,
        "maker_order_id": "ac928c66-ca53-498f-9c13-a110027a60e8",
        "taker_order_id": "132fb6ae-456b-4654-b4e0-d681ac05cea1",
        "time": "2014-11-07T08:19:27.028459Z",
        "product_id": "BTC-USD",
        "size": "5.23512",
        "price": "400.23",
        "side": "sell"
      }*/

      const dataList = this.kLineChart.getDataList()
      let lastData = dataList[dataList.length - 1]
      let msgTime = Date.parse(message.time)

      const kLineModel = {
        open: message.price,
        low: message.price,
        high: message.price,
        close: message.price,
        volume: parseFloat(message.size),
        timestamp: msgTime
      }

      // if the new message is within the same timeframe (period)
      // then we must apply this new data with the last timestamp
      // to keep the chart on the same candle
      let isP = this.isPeriod(lastData.timestamp, msgTime, this.category)
      if (isP) {
        kLineModel.open = lastData.open
        if (message.price <= lastData.low) {
          kLineModel.low = message.price
        }
        if (message.price >= lastData.high) {
          kLineModel.high = message.price
        }
        kLineModel.volume = (lastData.volume + parseFloat(message.size))
        kLineModel.timestamp = lastData.timestamp
        kLineModel.turnover = (kLineModel.open + kLineModel.close + kLineModel.high + kLineModel.low) / 4 * kLineModel.volume
        this.kLineChart.updateData(kLineModel)

      } else {
        kLineModel.turnover = (kLineModel.open + kLineModel.close + kLineModel.high + kLineModel.low) / 4 * kLineModel.volume
        this.kLineChart.updateData(kLineModel)
        /*dataList.unshift(kLineModel)
        this.kLineChart.applyMoreData(dataList)*/
      }



    },
    handleGetProductBook(message) {
      /*
      Response object format
      {
        "type": "getProductBook",
        "book": {
          "asks": [
            {
              "price": "38801.69",
              "size": "2.8877351",
              "num_orders": 1
            },
          ]
          "bids": [
            {
              "price": "38360.08",
              "size": "0.0003",
              "num_orders": 1
            },
          ]
          "sequence": 541837230
        },
      }*/
      message.book.bids.forEach(bid => {
        this.bids.push([bid.price, bid.size]);
      });
      message.book.asks.forEach(ask => {
        this.asks.push([ask.price, ask.size]);
      });
      this.bookInitialized = true
      this.calculate_midpoint()
    },
    handleSnapshot(message) {
      this.bids = []
      this.asks = []
      message.bids.forEach(bid => {
        this.bids.push([bid[0], bid[1]]);
      });
      message.asks.forEach(ask => {
        this.asks.push([ask[0], ask[1]]);
      });
      this.sortHighestToLowest(this.bids)
      this.sortLowestToHighest(this.asks)
      this.bids.length = this.bookArrayLimit
      this.asks.length = this.bookArrayLimit
      this.bookInitialized = true
      this.calculate_midpoint()
    },
    handleGetProductCandles(message) {
      this.kLineChart.clearData()

      const dataList = []
      message.candles.forEach(candle => {
        const kLineModel = {
          open: candle.open,
          low: candle.low,
          high: candle.high,
          close: candle.close,
          volume: candle.volume,
          timestamp: Date.parse(candle.time)
        }
        kLineModel.turnover = (kLineModel.open + kLineModel.close + kLineModel.high + kLineModel.low) / 4 * kLineModel.volume
        dataList.unshift(kLineModel)
      })


      this.kLineChart.applyNewData(dataList)
    },
    handleGetProducts(message) {
      /*
      {
      "base_currency":"LINK",
      "base_increment":0,
      "base_max_size":0,
      "base_min_size":0,
      "display_name": "LINK/USDC",
      "id":"LINK-USDC",
      "max_market_funds":0,
      "min_market_funds":0,
      "quote_currency":"USDC",
      "quote_increment":0,
      "status":"online",
      "status_message":"",
      "cancel_only":false,
      "limit_only":false,
      "post_only":false,
      "trading_disabled":false
      } */
      this.products = message.products
      this.productInputDisabled = false
    },
    handleGetExchanges(message) {
      this.exchanges = message.exchanges
    },

    /* Orderbook functions */
    calculate_midpoint() {
      if (this.bookInitialized) {
        let latestBidArr = this.bids.slice(0)[0]
        let latestAskArr = this.asks.slice(0)[0]
        let latestBid = parseFloat(latestBidArr[0])
        let latestAsk = parseFloat(latestAskArr[0])
        this.midpoint = ((latestBid + latestAsk) / 2).toFixed(3)
      }
    },
    update_book(data, side) {
      if (!this.bookInitialized) {
        return;
      }

      //todo: yeah this obv wont work for every currency type but its adding two rows for '1111' and '1111.00'
      //todo: switch on selected currency pair maybe and then round accordingly???
      //todo: better yet add a way to aggregate by price
      let price = parseFloat(data[1]).toFixed(4);
      let amount = parseFloat(data[2]);
      let zeroString = parseFloat("0.00000000")
      let index = -1

      if (side === "buy") {
        index = this.bids.findIndex(element => element[0] === price);
      } else {
        index = this.asks.findIndex(element => element[0] === price);
      }

      if (index > -1) {                 // If the index matches a price in the list then it is an update message
        if (amount === zeroString) {          // If the amount is zero
          this.removeBookOrder(side, index)   // remove the order
        } else {                              // otherwise update matching position in the book
          this.updateBookOrder(side, index, price, amount)
        }
      } else {                                // If the index is -1 then it is a new price that came in
        if (amount !== 0) {                   // If the amount is not zero insert new price
          this.insertBookOrder(side, price, amount)
        }
      }
      this.calculate_midpoint()
    },
    insertBookOrder(side, price, amount) {
      if (side === "buy") {
        if (price > this.bids[0][0]) {
          this.bids.unshift([price, amount]);
          if (this.bids.length >= this.bookArrayLimit) {
            this.bids.pop()
          }
        }
      } else {
        if (price < this.asks[0][0]) {
          this.asks.unshift([price, amount]);
          if (this.bids.length >= this.bookArrayLimit) {
            this.asks.pop()
          }
        }
      }
    },
    updateBookOrder(side, index, price, amount) {
      if (side === "buy") {
        this.bids[index] = [price, amount];
      } else {
        this.asks[index] = [price, amount];
      }
    },
    removeBookOrder(side, index) {
      if (side === "buy") {
        this.bids.splice(index, 1)
      } else {
        this.asks.splice(index, 1)
      }
    },
    sortHighestToLowest(data) {
      return data.sort((x, y) => parseFloat(y[0]) - parseFloat(x[0]))
    },
    sortLowestToHighest(data) {
      return data.sort((x, y) => parseFloat(x[0]) - parseFloat(y[0]))
    },

    /* Misc. */
    startChartTimer() {

      // get seconds remaining until next period
      let time = new Date()
      let secondsRemaining = (60 - time.getSeconds()) * 1000 - time.getMilliseconds();

      this.chartPeriodTimerFunc = setInterval(() => {
        const dataList = this.kLineChart.getDataList()
        const lastData = dataList[dataList.length - 1]
        const kLineModel = {
          open: lastData.close,
          low: lastData.close,
          high: lastData.close,
          close: lastData.close,
          volume: 0,
          timestamp: Date.now()
        }
        kLineModel.turnover = (kLineModel.open + kLineModel.close + kLineModel.high + kLineModel.low) / 4 * kLineModel.volume
        this.kLineChart.updateData(kLineModel)
        setInterval(this.chartPeriodTimerFunc, this.chartPeriod);
      }, secondsRemaining);
    },
    stopChartTimer() {
      clearInterval(this.chartPeriodTimerFunc);
    },
    isPeriod(oldDate, newDate, period) {
      let res = false
      const seconds = Math.floor((newDate - oldDate) / 1000);
      let interval = seconds / 60;

      switch (period) {
        case "1m":
          res = interval <= 1;
          break;
        case "5m":
          interval = seconds / 300;
          res = interval <= 1;
          break;
        case "15m":
          interval = seconds / 900;
          res = interval <= 1;
          break;
        case "30m":
          interval = seconds / 1800;
          res = interval <= 1;
          break;
        case "1h":
          interval = seconds / 3600;
          res = interval <= 1;
          break;
        case "4h":
          interval = seconds / 14400;
          res = interval <= 1;
          break;
        case "1d":
          interval = seconds / 86400;
          res = interval <= 1;
          break;
        case "1w":
          interval = seconds / 604800;
          res = interval <= 1;
          break;
        default:

      }

      return res
    }


  },
  destroyed: function () {
    window.removeEventListener("resize", this.handleWindowResize);
    this.stopChartTimer()
    this.socketSendUnsubReq()
    dispose('chart')
  }
}
</script>
<style scoped>
.chart-container {
  display: flex;
  flex-direction: column;
  flex: 1;
  height: calc(100%);
  width: calc(75%);
  border-bottom: 1px solid var(--border-color) !important;
}
.chart-top-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  height: 38px;
  font-size: 14px;
  color: #717484;
  border-bottom: 1px solid var(--border-color);
  box-sizing: border-box;
}
.chart-content {
  display: flex;
  flex-direction: row;
  height: calc(80% - 38px);
}

#chart {
  display: block !important;
  width: calc(100%) !important;
  height: 100%;
}

.connection-status {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  height: 38px;
  min-width: 49px !important;
  box-sizing: border-box;
  border-right: 0 solid var(--border-color);
  padding: 0;
  cursor: pointer;
}
.connection-time {
  text-align: center !important;

  border-right: 1px solid var(--border-color) !important;
}
.connection-time:hover {
  cursor: default !important;
}
.connection-message {
  margin: 5px !important;
  min-width: 75px;
  font-size: 10px;
}


.chart-top-bar .period {
  padding: 2px 6px;
  transition: all .2s;
  border-radius: 2px;
  margin-right: 6px;
  cursor: pointer;
}
.chart-top-bar .icon-wrapper {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  height: 38px;
  box-sizing: border-box;
  border-left: 1px solid var(--border-color);
  border-right: 0 solid var(--border-color);
  padding: 0;
  cursor: pointer;
}
.chart-top-bar .icon-wrapper a {
  padding: 0 12px;
  height: 100%;
  text-decoration: none !important;
}
.chart-top-bar .icon-wrapper svg {
  width: 20px;
  height: 20px;
  fill: #717484;
}
.chart-top-bar .icon-wrapper span {
  display: inline-block;
  margin-left: 4px;
}


.chart-tools-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  box-sizing: border-box;
  height: 100%;
  min-width: 50px !important;
  border-right: 1px solid var(--border-color);
}
.chart-tools-bar .icon-wrapper {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  height: 36px;
  width: 36px;
  margin-top: 8px;
  cursor: pointer;
  border-radius: 4px;
}
.chart-tools-bar .icon-wrapper svg {
  width: 36px;
  height: 36px;
  fill: #717484;
}
.chart-tools-bar .icon-wrapper svg:hover {
  background-color: var(--background-color-secondary) !important;
}

.divider {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  margin-bottom: 0 !important;
  margin-top: 8px !important;
  padding: 0 !important;
  height: 1px !important;
  width: 100%;
  border-radius: 0;
  opacity: 1!important;
  color: var(--border-color) !important;
  background-color: var(--border-color) !important;
  border: none !important;
}
.list-group-item {
  border: none !important;
}
.list-group-item:hover {
  cursor: pointer !important;
}
.dropdown-menu {
  padding: 10px !important;
}


.accordion {
  min-width: 300px;
}
.accordion-item, .accordion-button {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
  border-color: var(--border-color) !important;
}
.accordion-button {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  padding: 10px !important;
  font-size: 1rem;
}
.accordion-input {
  color: var(--text-primary-color) !important;
  background-color: var(--background-color-primary) !important;
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  padding: 8px 0 8px 8px !important;
  font-size: 1rem;
  text-align: left;
  border-radius: 0;
  cursor: default !important;
}
.accordion-divider {
  border-bottom: 1px solid var(--border-color) !important;
}
.accordion-body {
  padding: 0 0 0 10px !important;
}
.item {
  cursor: pointer;
  line-height: 1.5;
}

.form-select {
  width: 50% !important;
  border: 1px solid var(--border-color) !important;
  background-color: var(--background-color-primary) !important;
  cursor: pointer !important;
}
.form-select:focus {
  border: 1px solid var(--border-color) !important;
}

button:focus, select:focus, input:focus {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
}
.form-range {
  width: 50% !important;
  background-color: var(--background-color-primary) !important;
}
.form-range::-webkit-slider-runnable-track {
  background: var(--slider-color-left) !important;
}
.form-range::-webkit-slider-thumb {
  background: var(--slider-color-left) !important;
}
.form-range::-moz-range-track {
  background: var(--slider-color-left) !important;
}
.form-range::-moz-range-thumb {
  background: var(--slider-color-left) !important;
}
.form-range::-webkit-slider-thumb {
  background: var(--slider-color-left) !important;
}
.form-range::-moz-range-thumb {
  background: var(--slider-color-left) !important;
}
.form-range::-ms-thumb {
  background: var(--slider-color-left) !important;
}
.form-check-input:checked {
  background-color: #0d6efd !important;
  border-color: #0d6efd !important;
  border: 1px solid var(--border-color) !important;
}
.form-check-input, .form-check-input:before, .form-check-input:after, .form-check-input:focus {
  border: 1px solid var(--border-color) !important;
}


.red {
  color: red;
}
.green {
  color: green;
}
.bold {
  font-weight: bold;
}
.selected {
  color: var(--text-primary-color) !important;
}

#chartPeriodSelect {
  border: none !important;
}
#chartPeriodSelect:hover {
  background-color: var(--background-color-secondary) !important;
}



</style>
