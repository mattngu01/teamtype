
<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    data() {
        return {
            typedText: "",
        }
    },
    methods: {
    },
    watch: {
        typedText(oldTypedText: string, newTypedText: string) {
            // based on assumption that this is called on each char / each change
            let diff : number = Math.abs(oldTypedText.length - newTypedText.length);
            if (oldTypedText.length < newTypedText.length) {
                console.log("+" + newTypedText.slice(-diff))
            } else {
                console.log("-" + oldTypedText.slice(-diff))
            }
        }
    }
})

</script>
<template>
    <div> Hello world </div>
    <input for="test" spellcheck="false" autocapitalize="off" autocomplete="off" v-model="typedText" placeholder="Type here">
</template>